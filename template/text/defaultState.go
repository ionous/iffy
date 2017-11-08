package text

import (
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/template/chart"
)

type DefaultState struct {
	Engine
	PrevState
}

func (b DefaultState) next(c spec.Block, d chart.Directive) (ret DirectiveState, err error) {
	return b.advance(b, c, d)
}
