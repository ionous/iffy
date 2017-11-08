package text

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/kindOf"

	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

type Engine struct {
	Commander
	cmds Commands
	prev PrevStates
}

// StartJoin of one or more strings.
// Always returns true; the result exists to improve readability.
func (eng *Engine) span() (err error) {
	if cmd, e := eng.CreateCommand("join"); e != nil {
		err = e
	} else {
		err = eng.cmds.begin(cmd)
	}
	return
}

// because we are mixing text and evals, we expect the whole thing winds up being text. ( otherwise: what would we do with the intervening text. )
func (eng *Engine) convert(dirs []template.Directive) (ret *ops.Command, err error) {
	if len(eng.cmds.list) != 0 {
		err = errutil.New("engine reused")
	} else if cmd, e := eng.CreateCommand("join"); e != nil {
		err = e
	} else {
		eng.cmds.list = []*ops.Command{cmd}
		var state DirectiveState = DefaultState{Engine: eng}
		for _, d := range dirs {
			if n, e := state.next(d); e != nil {
				err = e
				break
			} else {
				state = n
			}
		}
		if err == nil {
			if cnt := len(eng.cmds.list); cnt != 1 {
				err = errutil.New("cmds error")
			} else {
				ret = eng.cmds.list[0]
			}
		}
	}
	return
}

func (eng *Engine) advance(q DirectiveState, d template.Directive) (ret DirectiveState, err error) {
	switch key, xs := d.Key, d.Expression; key {
	case "once":
		key = "stopping"
		fallthrough
	case "cycle", "shuffle":
		if xs != nil {
			err = errutil.New(key, "expected empty expression", xs)
		} else if seq, e := eng.newSequence(q, key); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = seq
		}
	case "if", "unless":
		if len(xs) == 0 {
			err = errutil.New("expected conditional")
		} else if cnd, e := eng.newCondition(q, xs, key != "if"); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = cnd
		}
	case "":
		if cmd, e := eng.CreateExpression(xs, kindOf.TypeTextEval); e != nil {
			err = e
		} else if e := eng.cmds.position(cmd); e != nil {
			err = e
		} else {
			ret = q // keep going in the same state
		}
	default:
		err = errutil.New("unknown key", key)
	}
	return
}
func (eng *Engine) newEnd(q DirectiveState) (ret EndState, err error) {
	if e := eng.span(); e != nil {
		err = e
	} else {
		eng.prev.push(q)
		ret = EndState{eng, 1}
	}
	return
}

func (eng *Engine) newSequence(q DirectiveState, group string) (ret SequenceState, err error) {
	if counter, e := eng.CreateName(group + " counter"); e != nil {
		err = e
	} else if count, e := eng.CreateCommand(group + " text"); e != nil {
		err = e
	} else if e := count.Position(counter); e != nil {
		err = e
	} else if e := eng.cmds.begin(count); e != nil {
		err = e
	} else if e := eng.span(); e != nil {
		err = e
	} else {
		eng.prev.push(q)
		ret = SequenceState{eng, 1}
	}
	return
}

func (eng *Engine) newCondition(q DirectiveState, xs postfix.Expression, invert bool) (ret ConditionState, err error) {
	if cmd, e := eng.CreateCommand("choose text"); e != nil {
		err = e
	} else if test, e := eng.CreateExpression(xs, kindOf.TypeBoolEval); e != nil {
		err = e
	} else if test, e := eng.invert(test, invert); e != nil {
		err = e
	} else if e := cmd.Position(test); e != nil {
		err = e
	} else if e := eng.cmds.begin(cmd); e != nil {
		err = e
	} else if e := eng.span(); e != nil {
		err = e
	} else {
		eng.prev.push(q)
		ret = ConditionState{eng, 1}
	}
	return
}

func (eng *Engine) invert(test *ops.Command, invert bool) (ret *ops.Command, err error) {
	if !invert {
		ret = test
	} else if inv, e := eng.CreateCommand("is not"); e != nil {
		err = e
	} else if e := inv.Position(test); e != nil {
		err = e
	} else {
		ret = inv
	}
	return
}
