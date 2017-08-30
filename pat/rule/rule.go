package rule

import (
	"github.com/ionous/iffy/ref/unique"
)

// Mandate defines an interface for creating rules.
type Mandate interface {
	Mandate(unique.Types, Rules) error
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

func (p Mandates) Mandate(patterns unique.Types, rules Rules) (err error) {
	for _, el := range p {
		if e := el.Mandate(patterns, rules); e != nil {
			err = e
			break
		}
	}
	return
}
