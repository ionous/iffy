package obj

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/coerce"
	"github.com/ionous/iffy/ref/enum"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/kind"
	r "reflect"
)

// UnpackValue returns the value of a property from an object, translating ids to objects.
func UnpackValue(or ObjectMap, obj rt.Object, name string, pv interface{}) (err error) {
	if pdst := r.ValueOf(pv); pdst.Kind() != r.Ptr {
		err = errutil.New(obj, name, "expected pointer outvalue", pdst.Type())
	} else if p, ok := obj.Property(name); !ok {
		err = errutil.New(obj, name, "unknown property")
	} else {
		cc := copyContext{or}
		dst := pdst.Elem()
		src := r.ValueOf(p.Value())
		if e := cc.copy(dst, src); e != nil {
			err = errutil.New(obj, name, "cant unpack", dst.Type(), "from", src.Type(), "because", e)
		}
	}
	return
}

// PackValue sets a property in a object to a value, translating objects to ids.
func PackValue(or ObjectMap, obj rt.Object, name string, v interface{}) (err error) {
	if p, ok := obj.Property(name); !ok {
		err = errutil.New(obj, name, "unknown property")
	} else {
		cc := copyContext{or}
		dst := r.New(p.Type()).Elem() // create a new destination for the value.
		src := r.ValueOf(v)
		if e := cc.copy(dst, src); e != nil {
			err = errutil.New(obj, name, "cant pack", dst.Type(), "from", src.Type(), "because", e)
		} else {
			err = p.SetValue(dst.Interface())
		}
	}
	return
}

type copyFun func(c *copyContext, dst, src r.Value) error

type copyContext struct {
	objs ObjectMap
}

func (cc *copyContext) copy(dst, src r.Value) (err error) {
	if ds, ss := dst.Kind() == r.Slice, src.Kind() == r.Slice; ds != ss {
		err = errutil.New("slice mismatch")
	} else {
		dt, st := dst.Type(), src.Type()
		if dt == st {
			dst.Set(src)
		} else if ds /*&& ss*/ {
			if cfn := getCopyFun(dt.Elem(), st.Elem()); cfn != nil {
				err = coerce.Slice(dst, src, func(dst, src r.Value) error {
					return cfn(cc, dst, src)
				})
			} else {
				err = coerce.Value(dst, src)
			}
		} else /*if !ds && !ss */ {
			if cfn := getCopyFun(dt, st); cfn != nil {
				err = cfn(cc, dst, src)
			} else {
				err = coerce.Value(dst, src)
			}
		}
	}
	return
}

func getCopyFun(dst, src r.Type) (ret copyFun) {
	switch {
	case dst.Kind() == r.Int && src.Kind() == r.String:
		ret = intFromChoice
	case src.Kind() == r.Int && dst.Kind() == r.String:
		ret = choiceFromInt

		// dst ident.Id src obj.RefObject because type mismatch"

	case src == kind.IdentId():
		ret = objFromId
	case dst == kind.IdentId():
		ret = idFromObj
	}
	return
}

func intFromChoice(c *copyContext, dst, src r.Value) (err error) {
	if !enum.Pack(dst, src) {
		err = errutil.New("couldnt pack enum")
	}
	return
}

func choiceFromInt(c *copyContext, dst, src r.Value) (err error) {
	if !enum.Unpack(dst, src) {
		err = errutil.New("couldnt pack enum")
	}
	return
}

func objFromId(c *copyContext, dst, src r.Value) (err error) {
	id := src.Interface().(ident.Id)
	if obj, ok := c.objs[id]; !ok {
		err = errutil.New("unknown object", id)
	} else {
		// seems to be trying to set the wrong way round
		// obj to id
		dst.Set(r.ValueOf(obj))
	}
	return
}

func idFromObj(c *copyContext, dst, src r.Value) (err error) {
	var id ident.Id
	if obj, ok := src.Interface().(rt.Object); !ok {
		err = errutil.New("src is not an object", src.Type())
	} else if obj != nil {
		id = obj.Id()
	}
	dst.Set(r.ValueOf(id))
	return
}
