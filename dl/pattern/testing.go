package pattern

import "github.com/ionous/iffy/dl/core"

func NewNamedParams(name string, from core.Assignment) *core.Arguments {
	return &core.Arguments{[]*core.Argument{{
		name, from,
	}}}
}
