package template

import (
	"github.com/ionous/iffy/spec"
)

type sequence struct {
	tcontext
	tprev
	tdepth
}

func (q sequence) advance(c spec.Block, t Token) (ret tstate, err error) {
	if plain(c, t) {
		ret = q
	} else {
		switch op := t.Str; op {
		case "or":
			endJoin(c)
			startJoin(c) // start a new join for the new section
			ret = q
		case "end":
			if prev, e := q.pop(); e != nil {
				err = e
			} else {
				ret = prev
				q.rollup(c) // end the array and sequence
			}
		default:
			ret, err = q.defaultAdvance(q, c, t)
		}
	}
	return
}

// end the ambient join, without changing state:
func endJoin(c spec.Block) {
	c.End()
}

func startJoin(c spec.Block) bool {
	return c.Cmd("join").Begin()
}
