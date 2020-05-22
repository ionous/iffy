package chart

import (
	"github.com/ionous/iffy/template/postfix"
)

// OperandState reads a single number, reference, or quote.
type OperandState interface {
	State
	GetOperand() (postfix.Function, error)
}

type OperandParser struct {
	OperandState
}

func (p *OperandParser) StateName() string {
	return "operand parser"
}

func (p *OperandParser) NewRune(r rune) (ret State) {
	var op OperandState
	switch {
	case isQuote(r):
		op = &QuoteParser{}
	case isDot(r):
		op = &FieldParser{}
	case isNumber(r) || r == '+' || r == '-':
		op = &NumParser{}
	case isLetter(r):
		op = &BooleanParser{}
	default:
		op = &EmptyOperand{r}
	}
	p.OperandState = op
	return op.NewRune(r)
}

func (p *OperandParser) GetExpression() (ret postfix.Expression, err error) {
	if op, e := p.OperandState.GetOperand(); e != nil {
		err = e
	} else if op != nil {
		ret = postfix.Expression{op}
	}
	return
}
