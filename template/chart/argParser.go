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

// GetArguments returns both the expression, and the number of separated arguments.
func (p ArgParser) GetArguments() (postfix.Expression, int, error) {
	return p.out, p.cnt, p.err
}

// NewRune starts with the first character of an argument;
// each arg is parsed via ArgParser.
func (p *ArgParser) NewRune(r rune) State {
	var sub ExpressionState
	if f := p.factory; f != nil {
		sub = p.factory.NewExpressionState()
	} else {
		sub = new(OperandParser)
	}
	return ParseChain(r, sub, Statement(func(r rune) (ret State) {
		if exp, e := sub.GetExpression(); e != nil {
			p.err = e
		} else if len(exp) > 0 {
			// doesn't shunt: operands are always left-to-right in any type of expression.
			p.cnt, p.out = p.cnt+1, append(p.out, exp...)
			if isSpace(r) {
				ret = MakeChain(spaces, p) // loop...
			}
		}
		return
	}))
}
