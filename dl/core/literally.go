package core

import (
	"github.com/ionous/iffy/rt/kind"
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
//
func literally(v interface{}, dstType r.Type) (ret interface{}, okay bool) {
	switch v := v.(type) {
	case bool:
		ret, okay = &Bool{v}, true
	case float64:
		ret, okay = &Num{v}, true
	case []float64:
		ret, okay = &Numbers{v}, true
	// -- string for a command.
	case string:
		// could be text or object --
		switch dstType {
		case kind.TextEval():
			ret, okay = &Text{v}, true
		case kind.ObjectEval():
			ret, okay = &Object{v}, true
		}
	case []string:
		switch dstType {
		case kind.TextListEval():
			ret, okay = &Texts{v}, true
		case kind.ObjListEval():
			ret, okay = &Objects{v}, true
		}
	default:
		{
			v := r.ValueOf(v)
			if kind.IsNumber(v.Kind()) {
				v := v.Convert(kind.Number()).Float()
				ret, okay = &Num{v}, true
			}
		}
	}
	return
}
