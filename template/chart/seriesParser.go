package chart

import (
	"github.com/ionous/iffy/template/postfix"
)

// SeriesParser reads a sequence of operand and operator phrases.
type SeriesParser struct {
	err error
	out postfix.Shunt
}

// NewRune starts on the first character of an operand or opening sub-phrase.
func (p *SeriesParser) NewRune(r rune) State {
	return p.operand(r)
}

func (p *SeriesParser) GetExpression() (ret postfix.Expression, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		ret, err = p.out.GetExpression()
	}
	return
}

// at the start of every operand we might have some opening paren, a bracket,
// or just some operand.
func (p *SeriesParser) operand(r rune) (ret State) {
	switch {
	case isOpenParen(r):
		p.out.BeginSubExpression()
		ret = MakeChain(spaces, Statement(p.operand))
	default:
		var a SubdirParser
		ret = ParseChain(r, &a, Statement(func(r rune) (ret State) {
			if exp, e := a.GetExpression(); e != nil {
				p.err = e
			} else if len(exp) > 0 {
				p.out.AddExpression(exp)
				ret = ParseChain(r, spaces, Statement(p.operator))
			}
			return
		}))
	}
	return
}

// after every argument can come operators or close parens or the end
// start on the first character of the operator.
// a pipe floats upward.
func (p *SeriesParser) operator(r rune) (ret State) {
	var b OperatorParser
	return ParseChain(r, &b, Statement(func(r rune) (ret State) {
		switch {
		case isCloseParen(r):
			p.out.EndSubExpression()
			ret = MakeChain(spaces, Statement(p.operator))
		default:
			ret = ParseChain(r, &b, Statement(func(r rune) (ret State) {
				if op, ok := b.GetOperator(); ok {
					p.out.AddFunction(op)
					ret = ParseChain(r, spaces, Statement(p.operand))
				}
				return
			}))
		}
		return
	}))
}

// looks for the end of sub-expressions before handling the next state.
func (p *SeriesParser) closing(r rune, next State) (ret State) {
	switch {
	case isCloseParen(r):
		p.out.EndSubExpression()
		ret = MakeChain(spaces, next)
	default:
		ret = next.NewRune(r)
	}
	return
}
