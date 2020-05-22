package chart

import "github.com/ionous/iffy/template/postfix"

// implements OperandState by returning nothing
type EmptyOperand struct{ r rune }

func (p *EmptyOperand) StateName() string {
	return "empty operand"
}

func (p *EmptyOperand) NewRune(r rune) State {
	return nil
}

func (p *EmptyOperand) GetOperand() (postfix.Function, error) {
	return nil, nil
}
