package template

import (
	"github.com/ionous/iffy/spec"
)

// Template contains tokens alternating between plain text and text which came from inside braces.
type Template struct {
	tokens    []Token
	gen       GenerateId
	parseExpr ExpressionParser
}

func (t Template) Convert(c spec.Block) (err error) {
	if len(t.tokens) == 1 {
		err = t.parseExpr(c, t.tokens[0].Str)
	} else {
		err = t.convertMulti(c)
	}
	return
}

// convert multiple evaluations.
// by definition the evaluations separated by plain text.
// ( if there were no text seperators, it would be one evaluation )

func (t Template) convertMulti(c spec.Block) (err error) {
	// because we are mixing text and evals, we expect the whole thing winds up being text. ( otherwise: what would we do with the intervening text. )
	// start a new join for the new section
	if c.Cmd("join").Begin() {
		if c.Cmds().Begin() {
			ctx := tcontext{t.gen, t.parseExpr}
			prev := tprev{}
			var state tstate = base{ctx, prev}
			for _, token := range t.tokens {
				if n, e := state.advance(c, token); e != nil {
					err = e
					break
				} else if n == nil {
					panic("state is nil")
				} else {
					state = n
				}
			}
			c.End()
		}
		c.End()
	}
	return
}
