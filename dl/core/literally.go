package core

import (
	"github.com/ionous/iffy/ref/kindOf"
	r "reflect"
)

// Transform converts values specified by a scriptinto values usable by the runtime.
// For instance, a string into a request for an object; an int into a number eval.
// See also express.NewTransform which can create templates from author specified strings.
// FIX: this uses core, but should it be a part of core?
func Transform(src r.Value, hint r.Type) (ret r.Value, err error) {
	if v := literally(src, hint); v != nil {
		ret = r.ValueOf(v)
	} else {
		ret = src
	}
	return
}

// literally allows users to specify primitive values for some evals.
//
// c.Cmd("strings", sliceOf.String("one", "two", "three"))
// c.Value(sliceOf.String("one", "two", "three"))
//
// c.Cmd("get").Begin() { c.Cmd("object", "@") c.Value("text") }
// c.Cmd("get", "@", "text")
//
func literally(src r.Value, dstType r.Type) (ret interface{}) {
	switch srcType := src.Type(); {
	case kindOf.Bool(srcType):
		v := src.Bool()
		ret = &Bool{v}

	case kindOf.Int(srcType):
		v := src.Int()
		if kindOf.NumListEval(dstType) {
			ret = &Numbers{[]float64{float64(v)}}
		} else {
			ret = &Num{float64(v)}
		}

	case kindOf.Float(srcType):
		v := src.Float()
		if kindOf.NumListEval(dstType) {
			ret = &Numbers{[]float64{v}}
		} else {
			ret = &Num{v}
		}

	// -- string for a command.
	case srcType.Kind() == r.String:
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
		case kindOf.TextListEval(dstType):
			ret = &Texts{[]string{v}}

		case kindOf.ObjListEval(dstType):
			ret = &Objects{[]string{v}}
		}

	case srcType.Kind() == r.Slice && kindOf.Float(srcType.Elem()):
		v := src.Interface().([]float64)
		ret = &Numbers{v}

	case srcType.Kind() == r.Slice && kindOf.String(srcType.Elem()):
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
