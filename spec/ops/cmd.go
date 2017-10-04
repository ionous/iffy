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
	idx, tgt := c.index, c.target
	if cnt := tgt.NumField(); idx+1 > cnt {
		err = errutil.New("too many arguments", tgt, "expected", cnt)
	} else if dst := tgt.Field(idx); !dst.IsValid() {
		err = errutil.New("couldnt get field", tgt, idx)
	} else {
		// check for the last value is an array:
		var autoel bool
		src := r.ValueOf(arg)
		if idx+1 == cnt && dst.Kind() == r.Slice && src.Kind() != r.Slice {
			if cmd, ok := arg.(*Command); ok {
				src, autoel = cmd.target.Addr(), true
			} else {
				box := r.New(dst.Type().Elem()).Elem()
				if e := c.setValue(box, src); e != nil {
					err = e
				} else {
					src, autoel = box, true
				}
			}
		}
		if err == nil {
			if autoel {
				if slice, e := appendValue(dst, src); e != nil {
					err = e
				} else {
					dst.Set(slice)
				}
			} else {
				if e := c.setField(dst, arg); e != nil {
					err = errutil.New("couldnt set field", tgt, idx, e)
				} else {
					c.index = idx + 1
				}
			}
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
			slice := dst
			for _, c := range v.els {
				if next, e := appendValue(slice, c.target.Addr()); e != nil {
					err = e
					break
				} else {
					slice = next
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

// seems to be duplication here between Commands the struct, which has append, and this.
func appendValue(slice r.Value, src r.Value) (ret r.Value, err error) {
	elType := slice.Type().Elem()
	if srcType := src.Type(); !srcType.AssignableTo(elType) {
		err = errutil.Fmt("incompatible element from %v to %v", srcType, elType)
	} else {
		ret = r.Append(slice, src)
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

// helper for managing error
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
