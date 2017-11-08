package text

import (
	"github.com/ionous/errutil"

	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

type ConditionState struct {
	*Engine
	Depth
}

func (q ConditionState) next(d template.Directive) (ret DirectiveState, err error) {
	switch key := d.Key; key {
	case "end":
		ret, err = q.rollup(q.Engine)
	case "else", "otherwise":
		q.cmds.end() // endJoin(c)
		if res, e := q.newBranch(d.Expression, false); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = res
		}
	case "unless":
		q.cmds.end() // endJoin(c)
		if res, e := q.newBranch(d.Expression, true); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = res
		}
	default:
		ret, err = q.advance(q, d)
	}
	return
}

func (q ConditionState) newBranch(x postfix.Expression, invert bool) (ret DirectiveState, err error) {
	if len(x) == 0 {
		if next, e := q.newEnd(q); e != nil {
			err = e
		} else {
			next.Depth = q.Depth
			ret = next
		}
	} else {
		if next, e := q.newCondition(q, x, invert); e != nil {
			err = e
		} else {
			next.Depth = 1 + q.Depth
			ret = next
		}
	}
	return
}
