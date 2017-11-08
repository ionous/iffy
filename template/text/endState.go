package text

import (
	"github.com/ionous/iffy/template"
)

type EndState struct {
	*Engine
	Depth
}

// EndState searches for an ending keyword
func (q EndState) next(d template.Directive) (ret DirectiveState, err error) {
	switch key := d.Key; key {
	case "end":
		ret, err = q.rollup(q.Engine)
	default:
		ret, err = q.advance(q, d)
	}
	return
}
