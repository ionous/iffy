package text

import (
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/template/chart"
)

type SequenceState struct {
	Engine
	PrevState
	Depth
}

func (q SequenceState) next(c spec.Block, d chart.Directive) (ret DirectiveState, err error) {
	switch key := d.Key; key {
	case "or":
		endJoin(c)
		startJoin(c) // start a new join for the new section
		ret = q
	case "end":
		if prev, e := q.pop(); e != nil {
			err = e
		} else {
			ret = prev
			q.rollup(c) // end the array and SequenceState
		}
	default:
		ret, err = q.advance(q, c, d)
	}
	return
}

// EndJoin of the last StartJoin.
func endJoin(c spec.Block) {
	c.End()
}

// StartJoin of one or more strings.
// Always returns true; the result exists to improve readability.
func startJoin(c spec.Block) bool {
	return c.Cmd("join").Begin()
}
