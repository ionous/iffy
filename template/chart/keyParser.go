package chart

import (
	"github.com/ionous/iffy/template"
)

// KeyParser reads a key and its optional following expression.
type KeyParser struct {
	err   error
	runes Runes
	exp   ExpressionState
}

// NewRune starts on the first letter of the key.
func (p *KeyParser) NewRune(r rune) (ret State) {
	if isLetter(r) {
		ret = p.runes.Accept(r, p)
	} else if isSpace(r) /*|| isCloseBracket(r) || isTrim(r)*/ {
		ret = MakeChain(spaces, p.exp)
	}
	return
}

func (p *KeyParser) GetDirective() (ret template.Directive, err error) {
	if e := p.err; e != nil {
		err = e
	} else if exp, e := p.exp.GetExpression(); e != nil {
		err = e
	} else {
		ret.Key = p.runes.String()
		ret.Expression = exp
	}
	return
}
