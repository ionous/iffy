package chart

import (
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
	"unicode"
)

// EmptyFactory creates parsers which match empty test.
type EmptyFactory struct{}

// EmptyParser reads empty text.
type EmptyParser struct{}

// AnyFactory creates parsers which match any series of lowercase letters
type AnyFactory struct{}

// AnyParser reads letters.
type AnyParser struct{ runes Runes }

func (EmptyFactory) NewExpressionState() ExpressionState {
	return EmptyParser{}
}
func (EmptyParser) GetExpression() (ret postfix.Expression, err error) {
	return
}
func (p EmptyParser) NewRune(r rune) (ret State) {
	if isSpace(r) {
		ret = p
	}
	return
}

func (f *AnyFactory) NewExpressionState() ExpressionState {
	return &AnyParser{}
}

func (p AnyParser) GetExpression() (ret postfix.Expression, err error) {
	if s := p.runes.String(); len(s) > 0 {
		arg := template.Reference([]string{s})
		ret = append(ret, arg)
	}
	return
}

func (p *AnyParser) NewRune(r rune) (ret State) {
	if unicode.IsLower(r) {
		ret = p.runes.Accept(r, p)
	}
	return
}
