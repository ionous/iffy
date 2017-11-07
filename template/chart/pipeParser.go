package chart

import (
	"github.com/ionous/iffy/template/postfix"
)

// Pipe parser reads an expression followed by an optional series of functions.
// Expression | Function | ...
type PipeParser struct {
	err error
	exp postfix.Expression
}

// NewRune starts on the first character of an operand or opening sub-phrase.
func (p *PipeParser) NewRune(r rune) State {
	return p.next(r, &ExpressionParser{})
}

func (p PipeParser) GetExpression() (postfix.Expression, error) {
	return p.exp, p.err
}

func (p *PipeParser) next(r rune, exp ExpressionState) State {
	return ParseChain(r, exp, Statement(func(r rune) (ret State) {
		if exp, e := exp.GetExpression(); e != nil {
			p.err = e
		} else if len(exp) > 0 {
			switch {
			case isPipe(r):
				ret = Statement(func(r rune) State {
					return p.next(r, &CallParser{arity: 1, out: exp})
				})
			default:
				p.exp = exp
			}
		}
		return
	}))
}
