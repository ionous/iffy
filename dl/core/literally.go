package core

import (
	"github.com/ionous/iffy/ref/kindOf"
	r "reflect"
)

// move to a different package? a sub-package?
// check the imports i guess.
type Xform struct{}

// returns src if no error but couldnt convert.
func (ts Xform) TransformValue(val interface{}, hint r.Type) (ret interface{}, err error) {
	if x := literally(val, hint); x != nil {
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
func literally(v interface{}, dstType r.Type) (ret interface{}) {
	switch v := v.(type) {
	case bool:
		ret = &Bool{v}
	case float64:
		ret = &Num{v}
	case []float64:
		ret = &Numbers{v}
	// -- string for a command.
	case string:
		// could be text or object --
		switch {
		case kindOf.TextEval(dstType):
			ret = &Text{v}
		case kindOf.ObjectEval(dstType):
			if v == "@" {
				ret = &TopObject{}
			} else {
				ret = &Object{v}
			}
		}
	case []string:
		switch {
		case kindOf.TextListEval(dstType):
			ret = &Texts{v}
		case kindOf.ObjListEval(dstType):
			ret = &Objects{v}
		}
	default:
		{
			v := r.ValueOf(v)
			if kindOf.Float(v.Type()) {
				v := v.Float()
				ret = &Num{v}
			} else if kindOf.Int(v.Type()) {
				v := v.Int()
				ret = &Num{float64(v)}
			}
		}
	}
	return
}
