package reduce

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
)

type Engine struct {
	factory Factory
}

// because we are mixing text and evals, we expect the whole thing winds up being text. ( otherwise: what would we do with the intervening text. )
func (eng Engine) reduce(c spec.Block, dirs []chart.Directive) (err error) {
	if startJoin(c) {
		var state DirectiveState = DefaultState{Engine: eng}
		for _, d := range dirs {
			if n, e := state.next(c, d); e != nil {
				err = e
				break
			} else {
				state = n
			}
		}
		endJoin(c)
	}
	return
}

func (eng Engine) advance(p DirectiveState, c spec.Block, d chart.Directive) (ret DirectiveState, err error) {
	switch key, x := d.Key, d.Expression; key {
	case "once":
		key = "stopping"
		fallthrough
	case "cycle", "shuffle":
		if x != nil {
			err = errutil.New(key, "expected empty expression", x)
		} else if seq, e := eng.newSequence(p, c, key); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = seq
		}
	case "if", "unless":
		if len(x) == 0 {
			err = errutil.New("expected conditional")
		} else if cnd, e := eng.newCondition(p, c, x, key == "if"); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = cnd
		}
	case "":
		if e := eng.factory.CreateExpression(c, x, kindOf.TypeTextEval); e != nil {
			err = e
		} else {
			ret = p // keep going in the same state
		}
	default:
		err = errutil.New("unknown key", key)
	}
	return
}

func (eng Engine) newSequence(p DirectiveState, c spec.Block, n string) (ret SequenceState, err error) {
	if counter, e := eng.factory.CreateName(n + " counter"); e != nil {
		err = e
	} else if c.Cmd(n+" text", counter).Begin() {
		ret = SequenceState{eng, PrevState{p}, 1}
		startJoin(c)
	}
	return
}

func (eng Engine) newCondition(p DirectiveState, c spec.Block, x postfix.Expression, is bool) (ret ConditionState, err error) {
	// FIX: without a dst hint, we cant choose anything but text.
	if c.Cmd("choose text").Begin() {
		if is {
			err = eng.factory.CreateExpression(c, x, kindOf.TypeBoolEval)
		} else if c.Cmd("is not").Begin() {
			err = eng.factory.CreateExpression(c, x, kindOf.TypeBoolEval)
			c.End()
		}
		if err == nil {
			ret = ConditionState{eng, PrevState{p}, 1}
			startJoin(c)
		}
	}
	return
}
