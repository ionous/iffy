package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/coerce"
	r "reflect"
)

type Command struct {
	xform  Transform
	target Target // output object we are building
	index  int
}

func (c *Command) Target() r.Value {
	return c.target.Addr()
}

func (c *Command) Position(arg interface{}) (err error) {
	if cnt := c.target.NumField(); c.index >= cnt {
		err = errutil.New("too many arguments", c.target, "expected", cnt)
	} else {
		field := c.target.Field(c.index)
		if !field.IsValid() {
			err = errutil.New("couldnt get field", c.target, c.index)
		} else if e := c.setField(field, arg); e != nil {
			err = errutil.New("couldnt get field", c.target, c.index, e)
		} else {
			c.index++
		}
	}
	return
}

func (c *Command) Assign(key string, arg interface{}) (err error) {
	field := c.target.FieldByName(key)
	if !field.IsValid() {
		err = errutil.Fmt("couldnt find field %T %v %v", c.target, c.target.Type(), key)
	} else if e := c.setField(field, arg); e != nil {
		err = errutil.New("field", key, e)
	}
	return
}

// dst is the field we are setting; v the value specified in the command script.
func (c *Command) setField(dst r.Value, v interface{}) (err error) {
	switch v := v.(type) {
	case *Command:
		// all commands are interfaces are implemented with pointers
		targetPtr := v.target.Addr()
		if e := coerce.Value(dst, targetPtr); e != nil {
			err = errutil.New("couldnt assign command", e)
		}

	case *Commands:
		if kind, isArray := arrayKind(dst.Type()); !isArray || kind != r.Interface {
			if !isArray {
				err = errutil.Fmt("trying to set an array to %v", dst.Type())
			} else {
				err = errutil.New("trying to set commands to", kind)
			}
		} else {
			slice, elType := dst, dst.Type().Elem()
			for _, c := range v.els {
				// all commands are interfaces are implemented with pointers
				rvalue := c.target.Addr()
				if from := rvalue.Type(); !from.AssignableTo(elType) {
					err = errutil.Fmt("incompatible element type. from: %v to: %v", from, elType)
					break
				} else {
					slice = r.Append(slice, rvalue)
				}
			}
			dst.Set(slice)
		}
	default:
		src := r.ValueOf(v)
		if dst.Kind() != r.Slice || src.Kind() != r.Slice {
			err = c.setValue(dst, src)
		} else {
			err = coerce.Slice(dst, src, func(del, sel r.Value) error {
				return c.setValue(del, sel)
			})
		}
	}
	return
}

func (c *Command) setValue(dst r.Value, src r.Value) (err error) {
	if v, e := xform(c.xform, src, dst.Type()); e != nil {
		err = e
	} else if !v.IsValid() {
		err = errutil.New("transform is empty")
	} else if e := coerce.Value(dst, v); e != nil {
		err = errutil.New("couldnt assign value", e)
	}
	return
}

// helper for managing errror
func xform(x Transform, src r.Value, hint r.Type) (ret r.Value, err error) {
	// if the destintation slot in the command is an interface -- ie. another command.
	if hint.Kind() == r.Interface {
		ret, err = x.TransformValue(src, hint)
	} else {
		ret = src
	}
	return
}

func arrayKind(rtype r.Type) (ret r.Kind, isArray bool) {
	if k := rtype.Kind(); k != r.Slice {
		ret = k
	} else {
		isArray = true
		ret = rtype.Elem().Kind()
	}
	return
}
