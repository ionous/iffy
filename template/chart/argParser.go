package chart

import (
	"github.com/ionous/iffy/template/postfix"
)

// ArgParser reads arguments specified by a function call.
type ArgParser struct {
	out     postfix.Expression
	cnt     int
	err     error
	factory ExpressionStateFactory
}

// MakeArgParser using a factory so that tests can mock out recursion;
// Normally, arguments can be any operand -- or any directive.
func MakeArgParser(f ExpressionStateFactory) ArgParser {
	return ArgParser{factory: f}
}

func (p ArgParser) GetExpression() (postfix.Expression, error) {
	return p.out, p.err
}

func (p ArgParser) GetArgs() (postfix.Expression, int, error) {
	return p.out, p.cnt, p.err
}

// NewRune starts with the first character of an argument;
// each arg is parsed via ArgParser.
func (p *ArgParser) NewRune(r rune) State {
	sub := p.factory.NewExpressionState()
	return parseChain(r, sub, Statement(func(r rune) (ret State) {
		if exp, e := sub.GetExpression(); e != nil {
			p.err = e
		} else if len(exp) > 0 {
			p.cnt, p.out = p.cnt+1, append(p.out, exp...)
			if isSpace(r) {
				ret = makeChain(spaces, p) // loop...
			}
		}
		return
	}))
}
