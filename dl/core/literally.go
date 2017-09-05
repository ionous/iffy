package core

import (
	"github.com/ionous/iffy/ref/kind"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// move to a different package? a sub-package?
// check the imports i guess.
type Xform struct{}

// returns src if no error but couldnt convert.
func (ts Xform) TransformValue(val interface{}, hint r.Type) (ret interface{}, err error) {
	if x, ok := literally(val, hint); ok {
		ret = x
	} else {
		ret = val
	}
	return
}

// literally allows users to specify primitive values for some evals.
//
// c.Cmd("texts", sliceOf.String("one", "two", "three"))
// c.Value(sliceOf.String("one", "two", "three"))
//
// c.Cmd("get").Begin() { c.Cmd("object", "@") c.Value("text") }
// c.Cmd("get", "@", "text")
func literally(v interface{}, dstType r.Type) (ret interface{}, okay bool) {
	switch v := v.(type) {
	case bool:
		ret, okay = &Bool{v}, true
	case float64:
		ret, okay = &Num{v}, true
	case []float64:
		ret, okay = &Numbers{v}, true
	case string:
		// could be text or object --
		switch dstType {
		case textEval:
			ret, okay = &Text{v}, true
		case objEval:
			ret, okay = &Object{v}, true
		}
	case []string:
		switch dstType {
		case textListEval:
			ret, okay = &Texts{v}, true
		case objListEval:
			ret, okay = &Objects{v}, true
		}
	default:
		{
			v := r.ValueOf(v)
			if kind.IsNumber(v.Kind()) {
				v := v.Convert(rFloat64).Float()
				ret, okay = &Num{v}, true
			}
		}
	}
	return
}

var rFloat64 r.Type = r.TypeOf((*float64)(nil)).Elem()
var textEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var objEval = r.TypeOf((*rt.ObjectEval)(nil)).Elem()
var textListEval = r.TypeOf((*rt.TextListEval)(nil)).Elem()
var objListEval = r.TypeOf((*rt.ObjListEval)(nil)).Elem()
