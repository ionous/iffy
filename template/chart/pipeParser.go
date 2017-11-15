package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
)

// Pipe parser reads an expression followed by an optional series of functions.
// Expression | Function | ...
type PipeParser struct {
	err error
	xs  postfix.Expression
}

// NewRune starts on the first character of an operand or opening sub-phrase.
func (p *PipeParser) NewRune(r rune) State {
	var expParser ExpressionParser
	return ParseChain(r, &expParser, p.after(0, &expParser))
}

func (p PipeParser) GetExpression() (postfix.Expression, error) {
	return p.xs, p.err
}

// after generates a state which reads the results of the passed expression parser.
func (p *PipeParser) after(n int, expParser ExpressionState) State {
	return Statement(func(r rune) (ret State) {
		if xs, e := expParser.GetExpression(); e != nil {
			p.err = e
		} else if cnt := len(xs); n > 0 && cnt == n {
			p.err = errutil.New("pipe should be followed by a call")
		} else {
			switch {
			case isPipe(r):
				ret = Statement(func(r rune) State {
					// pass the existing expression into the call parser.
					call := CallParser{arity: 1, out: xs}
					return ParseChain(r, &call, p.after(cnt, &call))
				})
			default:
				p.xs = xs
			}
		}
		return
	})
}
