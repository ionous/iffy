package text

import (
	"github.com/ionous/errutil"

	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

type ConditionState struct {
	*Engine
	PrevState
	Depth
}

func (n ConditionState) next(d template.Directive) (ret DirectiveState, err error) {
	switch key := d.Key; key {
	case "end":
		if prev, e := n.pop(); e != nil {
			err = e
		} else {
			n.rollup(n.Engine)
			ret = prev
		}
	case "else", "otherwise":
		n.end() // endJoin(c)
		if res, e := n.newBranch(d.Expression, false); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = res
		}
	case "unless":
		n.end() // endJoin(c)
		if res, e := n.newBranch(d.Expression, true); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = res
		}
	default:
		ret, err = n.advance(n, d)
	}
	return
}

func (n ConditionState) newBranch(x postfix.Expression, invert bool) (ret DirectiveState, err error) {
	if len(x) == 0 {
		if e := n.span(); e != nil {
			err = e
		} else {
			ret = EndState{n.Engine, PrevState{n}, n.Depth}
		}
	} else if cnd, e := n.newCondition(n, x, invert); e != nil {
		err = e
	} else {
		cnd.Depth += n.Depth
		ret = cnd
	}
	return
}
