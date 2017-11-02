package template

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/spec"
	"strings"
)

type tstate interface {
	advance(spec.Block, Token) (tstate, error)
	pop() (tstate, error)
}

type tprev struct {
	prev tstate
}

type tcontext struct {
	gen            NewName
	parseDirective DirectiveParser
}

// Restore a previous state, once this one is exhausted.
func (t tprev) pop() (ret tstate, err error) {
	if t.prev == nil {
		err = errutil.New("too many ends!")
	} else {
		ret = t.prev
	}
	return
}

// convert multiple evaluations.
// by definition the evaluations separated by plain text.
// ( if there were no text seperators, it would be one evaluation )
func (ctx tcontext) convertMulti(c spec.Block, ts []Token) (err error) {
	// because we are mixing text and evals, we expect the whole thing winds up being text. ( otherwise: what would we do with the intervening text. )
	// start a new join for the new section
	if startJoin(c) {
		prev := tprev{}
		var state tstate = base{ctx, prev}
		for _, token := range ts {
			if n, e := state.advance(c, token); e != nil {
				err = e
				break
			} else if n == nil {
				panic("state is nil")
			} else {
				state = n
			}
		}
		endJoin(c)
	}
	return
}

func (ctx tcontext) defaultAdvance(p tstat7e, c spec.Block, t Token) (ret tstate, err error) {
	if plain(c, t) {
		// keep going in the same state:
		ret = p
	} else {
		switch op := t.Str; op {
		case "once":
			ret, err = ctx.sequence(p, c, "stopping")
		case "cycle":
			ret, err = ctx.sequence(p, c, "cycle")
		case "shuffle":
			ret, err = ctx.sequence(p, c, "shuffle")
		default:
			parts := strings.Fields(t.Str)
			switch op, rest := parts[0], parts[1:]; op {

			case "if":
				if e := ctx.condition(c, rest, true); e != nil {
					err = errutil.New(op, e, t.Pos)
				} else {
					ret = conditional{ctx, tprev{p}, 1}
					startJoin(c)
				}
			case "unless":
				if e := ctx.condition(c, rest, false); e != nil {
					err = errutil.New(op, e, t.Pos)
				} else {
					ret = conditional{ctx, tprev{p}, 1}
					startJoin(c)
				}
			default:
				if e := ctx.parseDirective(c, parts, kindOf.TypeTextEval); e != nil {
					err = e
				} else {
					ret = p // keep going in the same state
				}
			}
		}
	}
	return
}

func (ctx tcontext) sequence(p tstate, c spec.Block, n string) (ret tstate, err error) {
	if c.Cmd(n+" text", ctx.gen.NewName(n+" counter")).Begin() {
		ret = sequence{ctx, tprev{p}, 1}
		startJoin(c)
	}
	return
}

func (ctx tcontext) condition(c spec.Block, rest []string, is bool) (err error) {
	if len(rest) == 0 {
		err = errutil.New("expected conditional")
	} else /*if len(rest) > 1 {
		err = errutil.New("currently supports a single condition")
	} else */{
		// FIX: without a dst hint, we cant choose anything but text.
		if c.Cmd("choose text").Begin() {
			if is {
				err = ctx.parseDirective(c, rest, kindOf.TypeBoolEval)
			} else if c.Cmd("is not").Begin() {
				err = ctx.parseDirective(c, rest, kindOf.TypeBoolEval)
				c.End()
			}
		}
	}
	return
}

func plain(c spec.Block, t Token) (ret bool) {
	if t.Plain {
		c.Cmd("text", t.Str)
		ret = true
	}
	return
}
