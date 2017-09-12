package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/coerce"
	"github.com/ionous/iffy/ref/enum"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type packFun func(rtm *Rtm, dst, src r.Value) error

func (rtm *Rtm) Pack(dst, src r.Value) (err error) {
	if ds, ss := dst.Kind() == r.Slice, src.Kind() == r.Slice; ds != ss {
		err = errutil.New("slice mismatch")
	} else {
		dt, st := dst.Type(), src.Type()
		if dt == st {
			dst.Set(src)
		} else if ds /*&& ss*/ {
			if cfn := getCopyFun(dt.Elem(), st.Elem()); cfn != nil {
				err = coerce.Slice(dst, src, func(dst, src r.Value) error {
					return cfn(rtm, dst, src)
				})
			} else {
				err = coerce.Value(dst, src)
			}
		} else /*if !ds && !ss */ {
			if cfn := getCopyFun(dt, st); cfn != nil {
				err = cfn(rtm, dst, src)
			} else {
				err = coerce.Value(dst, src)
			}
		}
	}
	return
}

func getCopyFun(dst, src r.Type) (ret packFun) {
	switch {

	case kindOf.IdentId(dst):
		ret = idFromObj

	case kindOf.IdentId(src):
		ret = objFromId

	case dst.Kind() == r.Int && src.Kind() == r.String:
		ret = intFromChoice

	case dst.Kind() == r.String && src.Kind() == r.Int:
		ret = choiceFromInt

		// asking for an eval, presumably given for a primitive
	case dst.Kind() == r.Interface:
		switch {
		case kindOf.BoolEval(dst):
			ret = toBoolEval
		case kindOf.NumberEval(dst):
			ret = toNumberEval
		case kindOf.TextEval(dst):
			ret = toTextEval
		case kindOf.ObjectEval(dst):
			ret = toObjEval
		case kindOf.NumListEval(dst):
			ret = toNumListEval
		case kindOf.TextListEval(dst):
			ret = toTextListEval
		case kindOf.ObjListEval(dst):
			ret = toObjListEval
		}

		// given an eval, presumably asking for a primitive
	case src.Kind() == r.Interface:
		switch {
		case kindOf.BoolEval(src):
			ret = fromBoolEval
		case kindOf.NumberEval(src):
			ret = fromNumberEval
		case kindOf.TextEval(src):
			ret = fromTextEval
		case kindOf.ObjectEval(src):
			ret = fromObjEval
		case kindOf.NumListEval(src):
			ret = fromNumListEval
		case kindOf.TextListEval(src):
			ret = fromTextListEval
		case kindOf.ObjListEval(src):
			ret = fromObjListEval
		}
	}
	return
}

func intFromChoice(rtm *Rtm, dst, src r.Value) (err error) {
	if !enum.Pack(dst, src) {
		err = errutil.New("couldnt pack enum")
	}
	return
}
func choiceFromInt(rtm *Rtm, dst, src r.Value) (err error) {
	if !enum.Unpack(dst, src) {
		err = errutil.New("couldnt pack enum")
	}
	return
}
func objFromId(rtm *Rtm, dst, src r.Value) (err error) {
	id := src.Interface().(ident.Id)
	if obj, ok := rtm.Objects[id]; !ok {
		err = errutil.New("unknown object", id)
	} else {
		// seems to be trying to set the wrong way round
		// obj to id
		dst.Set(r.ValueOf(obj))
	}
	return
}
func idFromObj(rtm *Rtm, dst, src r.Value) (err error) {
	var id ident.Id
	if obj, ok := src.Interface().(rt.Object); !ok {
		err = errutil.New("src is not an object", src.Type())
	} else if obj != nil {
		id = obj.Id()
	}
	dst.Set(r.ValueOf(id))
	return
}
func fromBoolEval(rtm *Rtm, dst, src r.Value) (err error) {
	eval := src.Interface().(rt.BoolEval)
	if v, e := eval.GetBool(rtm); e != nil {
		err = e
	} else {
		err = coerce.Value(dst, r.ValueOf(v))
	}
	return
}
func fromNumberEval(rtm *Rtm, dst, src r.Value) (err error) {
	eval := src.Interface().(rt.NumberEval)
	if v, e := eval.GetNumber(rtm); e != nil {
		err = e
	} else {
		err = coerce.Value(dst, r.ValueOf(v))
	}
	return
}
func fromTextEval(rtm *Rtm, dst, src r.Value) (err error) {
	eval := src.Interface().(rt.TextEval)
	if v, e := eval.GetText(rtm); e != nil {
		err = e
	} else {
		err = coerce.Value(dst, r.ValueOf(v))
	}
	return
}
func fromObjEval(rtm *Rtm, dst, src r.Value) (err error) {
	eval := src.Interface().(rt.ObjectEval)
	if v, e := eval.GetObject(rtm); e != nil {
		err = e
	} else {
		// recurse since we dont know the dst type.
		err = rtm.Pack(dst, r.ValueOf(v))
	}
	return
}
func fromNumListEval(rtm *Rtm, dst, src r.Value) (err error) {
	eval := src.Interface().(rt.NumListEval)
	if v, e := eval.GetNumberStream(rtm); e != nil {
		err = e
	} else {
		err = coerce.Value(dst, r.ValueOf(v))
	}
	return
}
func fromTextListEval(rtm *Rtm, dst, src r.Value) (err error) {
	eval := src.Interface().(rt.TextListEval)
	if v, e := eval.GetTextStream(rtm); e != nil {
		err = e
	} else {
		err = coerce.Value(dst, r.ValueOf(v))
	}
	return
}
func fromObjListEval(rtm *Rtm, dst, src r.Value) (err error) {
	eval := src.Interface().(rt.ObjListEval)
	if v, e := eval.GetObjectStream(rtm); e != nil {
		err = e
	} else {
		err = coerce.Value(dst, r.ValueOf(v))
	}
	return
}

func toBoolEval(rtm *Rtm, dst, src r.Value) (err error) {
	var v bool
	if e := coerce.Value(r.ValueOf(&v).Elem(), src); e != nil {
		err = e
	} else {
		dst.Set(r.ValueOf(&core.Bool{v}))
	}
	return
}
func toNumberEval(rtm *Rtm, dst, src r.Value) (err error) {
	var v float64
	if e := coerce.Value(r.ValueOf(&v).Elem(), src); e != nil {
		err = e
	} else {
		dst.Set(r.ValueOf(&core.Num{v}))
	}
	return
}
func toTextEval(rtm *Rtm, dst, src r.Value) (err error) {
	var v string
	if e := coerce.Value(r.ValueOf(&v).Elem(), src); e != nil {
		err = e
	} else {
		dst.Set(r.ValueOf(&core.Text{v}))
	}
	return
}
func toObjEval(rtm *Rtm, dst, src r.Value) (err error) {
	var v ident.Id
	if e := rtm.Pack(r.ValueOf(&v).Elem(), src); e != nil {
		err = e
	} else {
		dst.Set(r.ValueOf(&core.Object{v.Name}))
	}
	return
}
func toNumListEval(rtm *Rtm, dst, src r.Value) (err error) {
	var vs []float64
	if e := coerce.Value(r.ValueOf(&vs).Elem(), src); e != nil {
		err = e
	} else {
		dst.Set(r.ValueOf(&core.Numbers{vs}))
	}
	return
}
func toTextListEval(rtm *Rtm, dst, src r.Value) (err error) {
	var vs []string
	if e := coerce.Value(r.ValueOf(&vs).Elem(), src); e != nil {
		err = e
	} else {
		dst.Set(r.ValueOf(&core.Texts{vs}))
	}
	return
}
func toObjListEval(rtm *Rtm, dst, src r.Value) (err error) {
	return errutil.New("not implemented")
}
