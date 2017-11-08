package text

import (
	"github.com/ionous/iffy/template"
)

type SequenceState struct {
	*Engine
	Depth
}

func (q SequenceState) next(d template.Directive) (ret DirectiveState, err error) {
	switch key := d.Key; key {
	case "or":
		q.cmds.end() // endJoin(c)
		q.span()     // start a new join for the new section
		ret = q
	case "end":
		ret, err = q.rollup(q.Engine) // end the array and SequenceState
	default:
		ret, err = q.advance(q, d)
	}
	return
}
