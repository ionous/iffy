package rule

import (
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/rt"
)

// Mandate defines an interface for creating runtime patterns.
type Mandate interface {
	Mandate(RuleFactory) error
}

type RuleFactory interface {
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
	*ContinueAfter
	*ContinueBefore
	*Determine
	*NumberRule
	*NumListRule
	*ObjectRule
	*ObjListRule
	*RunRule
	*TextListRule
	*TextRule
}

type Mandates []Mandate

func (p Mandates) Mandate(fac RuleFactory) (err error) {
	for _, el := range p {
		if e := el.Mandate(fac); e != nil {
			err = e
			break
		}
	}
	return
}
