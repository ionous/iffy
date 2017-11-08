package text

import (
	"github.com/ionous/iffy/template"
)

type DefaultState struct {
	*Engine
	PrevState
}

func (b DefaultState) next(d template.Directive) (ret DirectiveState, err error) {
	return b.advance(b, d)
}
