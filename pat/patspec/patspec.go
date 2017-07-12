package patspec

import (
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/rt"
)

// MakePattern defines an interface for creating patterns.
type MakePattern interface {
	// Generate aka finalize ( in the sense of an album or collection. )
	Generate(PatternFactory) error
}

type PatternFactory interface {
	AddBool(string, rt.BoolEval, rt.BoolEval) error
	AddNumber(string, rt.BoolEval, rt.NumberEval) error
	AddText(string, rt.BoolEval, rt.TextEval) error
	AddObject(string, rt.BoolEval, rt.ObjectEval) error
	AddNumList(string, rt.BoolEval, rt.NumListEval) error
	AddTextList(string, rt.BoolEval, rt.TextListEval) error
	AddObjList(string, rt.BoolEval, rt.ObjListEval) error
	AddExecList(string, rt.BoolEval, rt.Execute, pat.Flags) error
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

type PatternSpecs []MakePattern

func (p PatternSpecs) Generate(fac PatternFactory) (err error) {
	for _, el := range p {
		if e := el.Generate(fac); e != nil {
			err = e
			break
		}
	}
	return
}
