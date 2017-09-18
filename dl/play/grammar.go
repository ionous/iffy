package play

import (
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/parser"
)

type Grammar struct {
	Match parser.Scanner // should be all of / any of.
}

type Mandate struct {
	rules.Mandate
}

func (a *Grammar) Define(f *Facts) (nil error) {
	f.Grammar.Match = append(f.Grammar.Match, a.Match)
	return
}

func (a *Mandate) Define(f *Facts) (nil error) {
	f.Mandates = append(f.Mandates, a.Mandate)
	return
}
