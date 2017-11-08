package text

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
)

type ConditionState struct {
	Engine
	PrevState
	Depth
}

func (n ConditionState) next(c spec.Block, d chart.Directive) (ret DirectiveState, err error) {
	switch key := d.Key; key {
	case "end":
		if prev, e := n.pop(); e != nil {
			err = e
		} else {
			n.rollup(c)
			ret = prev
		}
	case "else", "otherwise":
		endJoin(c)
		if res, e := n.newBranch(c, d.Expression, true); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = res
		}
	case "unless":
		endJoin(c)
		if res, e := n.newBranch(c, d.Expression, false); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = res
		}
	default:
		ret, err = n.advance(n, c, d)
	}
	return
}

func (n ConditionState) newBranch(c spec.Block, x postfix.Expression, useIf bool) (ret DirectiveState, err error) {
	if len(x) == 0 {
		ret = EndState{n.Engine, PrevState{n}, n.Depth}
		startJoin(c)
	} else if cnd, e := n.newCondition(n, c, x, useIf); e != nil {
		err = e
	} else {
		cnd.Depth += n.Depth
		ret = cnd
	}
	return
}
