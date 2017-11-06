package chart

import (
	"github.com/ionous/iffy/template/postfix"
)

// Pipe parser reads an expression followed by an optional series of functions.
// Expression | Function | ...
type PipeParser struct {
	err error
	out postfix.Pipe
}

// NewRune starts on the first character of an operand or opening sub-phrase.
func (p *PipeParser) NewRune(r rune) State {
	return p.next(r, &ExpressionParser{})
}

func (p *PipeParser) pipe(r rune) (ret State) {
	return p.next(r, &CallParser{arity: 1})
}

func (p *PipeParser) GetExpression() (ret postfix.Expression, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		ret, err = p.out.GetExpression()
	}
	return
}

func (p *PipeParser) next(r rune, exp ExpressionState) State {
	return ParseChain(r, exp, Statement(func(r rune) (ret State) {
		if res, e := exp.GetExpression(); e != nil {
			p.err = e
		} else if res != nil {
			// add each element of the expression:
			for _, x := range res {
				p.out.AddFunction(x)
			}
			switch {
			case isPipe(r):
				if e := p.out.AddPipe(); e != nil {
					p.err = e
				} else {
					ret = Statement(p.pipe)
				}
			}
		}
		return
	}))
}
