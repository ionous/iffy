package chart

import (
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// OperandState reads a single number, reference, or quote.
type OperandState interface {
	State
	GetOperand() (postfix.Function, error)
}

type OperandParser struct {
	OperandState
}

func (p *OperandParser) NewRune(r rune) (ret State) {
	var op OperandState
	switch {
	case isQuote(r):
		op = &QuoteParser{}
	case isLetter(r):
		op = &FieldParser{}
	case isNumber(r) || r == '+' || r == '-':
		op = &NumParser{}
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

type EmptyOperand struct{ r rune }

func (p *EmptyOperand) NewRune(r rune) State { return nil }

func (p EmptyOperand) GetOperand() (postfix.Function, error) {
	return nil, nil
}

func (p QuoteParser) GetOperand() (ret postfix.Function, err error) {
	if r, e := p.GetString(); e != nil {
		err = e
	} else {
		ret = types.Quote(r)
	}
	return
}

func (p FieldParser) GetOperand() (ret postfix.Function, err error) {
	if r, e := p.GetFields(); e != nil {
		err = e
	} else if b, ok := boolean(r); ok {
		ret = b
	} else {
		ret = types.Reference(r)
	}
	return
}

func boolean(r []string) (ret postfix.Function, okay bool) {
	if len(r) == 1 {
		switch r[0] {
		case "true":
			ret, okay = types.Boolean(true), true
		case "false":
			ret, okay = types.Boolean(false), true
		}
	}
	return
}

func (p NumParser) GetOperand() (ret postfix.Function, err error) {
	if n, e := p.GetValue(); e != nil {
		err = e
	} else {
		ret = types.Number(n)
	}
	return
}
