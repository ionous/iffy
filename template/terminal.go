package template

import (
	"github.com/ionous/iffy/spec"
)

type terminal struct {
	tcontext
	tprev
	tdepth
}

func (q terminal) advance(c spec.Block, t Token) (ret tstate, err error) {
	if plain(c, t) {
		ret = q
	} else {
		switch op := t.Str; op {
		case "end":
			if prev, e := q.pop(); e != nil {
				err = e
			} else {
				ret = prev
				q.rollup(c)
			}
		default:
			ret, err = q.defaultAdvance(q, c, t)
		}
	}
	return
}

type tdepth int

func (d tdepth) rollup(c spec.Block) {
	endJoin(c)
	for i := 0; i < int(d); i++ {
		c.End()
	}
}
