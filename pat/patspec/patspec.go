package patspec

import (
	"github.com/ionous/iffy/rt"
)

type Pattern interface {
	Generate(PatternFactory) error
}

type PatternFactory interface {
	Bool(string, rt.BoolEval, rt.BoolEval) error
	Number(string, rt.BoolEval, rt.NumberEval) error
	Text(string, rt.BoolEval, rt.TextEval) error
	Object(string, rt.BoolEval, rt.ObjectEval) error
	NumList(string, rt.BoolEval, rt.NumListEval) error
	TextList(string, rt.BoolEval, rt.TextListEval) error
	ObjList(string, rt.BoolEval, rt.ObjListEval) error
}

type Commands struct {
	*BoolRule
	*NumberRule
	*TextRule
	*ObjectRule
	*NumListRule
	*TextListRule
	*ObjListRule
	*Determine
}
