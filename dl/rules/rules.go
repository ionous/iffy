package rules

import (
	"github.com/ionous/iffy/pat"
)

// Mandate defines an interface for creating rules.
type Mandate interface {
	Mandate(pat.Contract) error
}

type MandateCmds struct {
	*BoolRule
	*NumberRule
	*NumListRule
	*ObjectRule
	*ObjListRule
	*RunRule
	*TextListRule
	*TextRule
}

type TimingCmds struct {
	*ContinueAfter
	*ContinueBefore
}

type RuntimeCmds struct {
	*Determine
}

type Commands struct {
	MandateCmds
	TimingCmds
	RuntimeCmds
}

type Mandates []Mandate

func (p Mandates) Mandate(rules pat.Contract) (err error) {
	for _, el := range p {
		if e := el.Mandate(rules); e != nil {
			err = e
			break
		}
	}
	return
}
