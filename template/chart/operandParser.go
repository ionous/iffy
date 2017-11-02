package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
)

// OperandState reads a single operand.
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
		op = &InvalidOperand{}
	}
	p.OperandState = op
	return op.NewRune(r)
}

type InvalidOperand struct{}

func (p *InvalidOperand) NewRune(r rune) State { return nil }

func (p *InvalidOperand) GetOperand() (postfix.Function, error) {
	return nil, errutil.New("invalid operand")
}

func (p QuoteParser) GetOperand() (ret postfix.Function, err error) {
	if r, e := p.GetString(); e != nil {
		err = e
	} else {
		ret = Quote(r)
	}
	return
}

func (p FieldParser) GetOperand() (ret postfix.Function, err error) {
	if r, e := p.GetFields(); e != nil {
		err = e
	} else {
		ret = Reference(r)
	}
	return
}

func (p NumParser) GetOperand() (ret postfix.Function, err error) {
	if n, e := p.GetValue(); e != nil {
		err = e
	} else {
		ret = Number(n)
	}
	return
}
