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
	stack []*ops.Command
}

func (eng *Engine) end() {
	eng.stack = eng.stack[:len(eng.stack)-1]
}

func (eng *Engine) begin(cmd *ops.Command) (err error) {
	if e := eng.position(cmd); e != nil {
		err = e
	} else {
		eng.stack = append(eng.stack, cmd)
	}
	return
}

func (eng *Engine) position(cmd *ops.Command) (err error) {
	if cnt := len(eng.stack); cnt == 0 {
		err = errutil.New("stack underflow")
	} else {
		top := eng.stack[cnt-1]
		if e := top.Position(cmd); e != nil {
			err = e
		}
	}
	return
}

// StartJoin of one or more strings.
// Always returns true; the result exists to improve readability.
func (eng *Engine) span() (err error) {
	if cmd, e := eng.CreateCommand("join"); e != nil {
		err = e
	} else {
		err = eng.begin(cmd)
	}
	return
}

// because we are mixing text and evals, we expect the whole thing winds up being text. ( otherwise: what would we do with the intervening text. )
func (eng *Engine) convert(dirs []template.Directive) (ret *ops.Command, err error) {
	if len(eng.stack) != 0 {
		err = errutil.New("engine reused")
	} else if cmd, e := eng.CreateCommand("join"); e != nil {
		err = e
	} else {
		eng.stack = []*ops.Command{cmd}
		var state DirectiveState = DefaultState{Engine: eng}
		for _, d := range dirs {
			if n, e := state.next(d); e != nil {
				err = e
				break
			} else {
				state = n
			}
		}
		if cnt := len(eng.stack); cnt != 1 {
			err = errutil.New("stack error", eng.stack[1].Target().Type().String())
		} else {
			ret = eng.stack[0]
		}
	}
	return
}

func (eng *Engine) advance(p DirectiveState, d template.Directive) (ret DirectiveState, err error) {
	switch key, xs := d.Key, d.Expression; key {
	case "once":
		key = "stopping"
		fallthrough
	case "cycle", "shuffle":
		if xs != nil {
			err = errutil.New(key, "expected empty expression", xs)
		} else if seq, e := eng.newSequence(p, key); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = seq
		}
	case "if", "unless":
		if len(xs) == 0 {
			err = errutil.New("expected conditional")
		} else if cnd, e := eng.newCondition(p, xs, key != "if"); e != nil {
			err = errutil.New(key, e)
		} else {
			ret = cnd
		}
	case "":
		if cmd, e := eng.CreateExpression(xs, kindOf.TypeTextEval); e != nil {
			err = e
		} else if e := eng.position(cmd); e != nil {
			err = e
		} else {
			ret = p // keep going in the same state
		}
	default:
		err = errutil.New("unknown key", key)
	}
	return
}

func (eng *Engine) newSequence(p DirectiveState, n string) (ret SequenceState, err error) {
	if counter, e := eng.CreateName(n + " counter"); e != nil {
		err = e
	} else if count, e := eng.CreateCommand(n + " text"); e != nil {
		err = e
	} else if e := count.Position(counter); e != nil {
		err = e
	} else if e := eng.begin(count); e != nil {
		err = e
	} else if e := eng.span(); e != nil {
		err = e
	} else {
		ret = SequenceState{eng, PrevState{p}, 1}
	}
	return
}

func (eng *Engine) newCondition(p DirectiveState, xs postfix.Expression, invert bool) (ret ConditionState, err error) {
	if cmd, e := eng.CreateCommand("choose text"); e != nil {
		err = e
	} else if test, e := eng.CreateExpression(xs, kindOf.TypeBoolEval); e != nil {
		err = e
	} else if test, e := eng.invert(test, invert); e != nil {
		err = e
	} else if e := cmd.Position(test); e != nil {
		err = e
	} else if e := eng.begin(cmd); e != nil {
		err = e
	} else if e := eng.span(); e != nil {
		err = e
	} else {
		ret = ConditionState{eng, PrevState{p}, 1}
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
