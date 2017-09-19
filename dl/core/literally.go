package core

import (
	"github.com/ionous/iffy/ref/kindOf"
	r "reflect"
)

// move to a different package? a sub-package?
// check the imports i guess.
type Xform struct{}

// returns src if no error but couldnt convert.
func (ts Xform) TransformValue(src r.Value, hint r.Type) (ret r.Value, err error) {
	if v := literally(src, hint); v != nil {
		ret = r.ValueOf(v)
	} else {
		ret = src
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
func literally(src r.Value, dstType r.Type) (ret interface{}) {
	switch rtype := src.Type(); {
	case kindOf.Bool(rtype):
		v := src.Bool()
		ret = &Bool{v}

	case kindOf.Int(rtype):
		v := src.Int()
		ret = &Num{float64(v)}

	case kindOf.Float(rtype):
		v := src.Float()
		ret = &Num{v}

	// -- string for a command.
	case rtype.Kind() == r.String:
		// could be text or object --
		v := src.String()
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

	case rtype.Kind() == r.Slice && kindOf.Float(rtype.Elem()):
		v := src.Interface().([]float64)
		ret = &Numbers{v}

	case rtype.Kind() == r.Slice && kindOf.String(rtype.Elem()):
		v := src.Interface().([]string)
		switch {
		case kindOf.TextListEval(dstType):
			ret = &Texts{v}
		case kindOf.ObjListEval(dstType):
			ret = &Objects{v}
		}
	}
	return
}
