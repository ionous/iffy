package text

import (
	"github.com/ionous/iffy/template"
)

type EndState struct {
	*Engine
	PrevState
	Depth
}

func (q EndState) next(d template.Directive) (ret DirectiveState, err error) {
	switch key := d.Key; key {
	case "end":
		if prev, e := q.pop(); e != nil {
			err = e
		} else {
			ret = prev
			q.rollup(q.Engine)
		}
	default:
		ret, err = q.advance(q, d)
	}
	return
}

type Depth int

func (d Depth) rollup(eng *Engine) {
	eng.end() // end span
	for i := 0; i < int(d); i++ {
		eng.end()
	}
}
