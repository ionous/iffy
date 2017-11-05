package chart

import (
	"github.com/ionous/iffy/template/postfix"
	"unicode"
)

// creates parsers which match empty test
type EmptyFactory struct{}

// creates directives
type EmptyParser struct{}

// AnyFactory creates parsers which match any series of lowercase letters
type AnyFactory struct{}

// creates directives
type AnyParser struct{ runes Runes }

// NewExpressionState
func (EmptyFactory) NewExpressionState() ExpressionState           { return EmptyParser{} }
func (EmptyParser) NewRune(rune) State                             { return nil }
func (EmptyParser) GetExpression() (x postfix.Expression, e error) { return }

// NewExpressionState
func (f *AnyFactory) NewExpressionState() ExpressionState {
	return &AnyParser{}
}

func (p AnyParser) GetExpression() (ret postfix.Expression, err error) {
	if s := p.runes.String(); len(s) > 0 {
		arg := Reference([]string{s})
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
