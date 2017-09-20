package template

import (
	"github.com/ionous/iffy/spec"
)

type base struct {
	tcontext
	tprev
}

func (b base) advance(c spec.Block, t Token) (ret tstate, err error) {
	return b.defaultAdvance(b, c, t)
}
