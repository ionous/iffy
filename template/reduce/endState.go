package reduce

import (
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/template/chart"
)

type EndState struct {
	Engine
	PrevState
	Depth
}

func (q EndState) next(c spec.Block, d chart.Directive) (ret DirectiveState, err error) {
	switch key := d.Key; key {
	case "end":
		if prev, e := q.pop(); e != nil {
			err = e
		} else {
			ret = prev
			q.rollup(c)
		}
	default:
		ret, err = q.advance(q, c, d)
	}
	return
}

type Depth int

func (d Depth) rollup(c spec.Block) {
	endJoin(c)
	for i := 0; i < int(d); i++ {
		c.End()
	}
}
