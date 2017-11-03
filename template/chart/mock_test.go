package chart

import (
	"github.com/ionous/iffy/template/postfix"
	"unicode"
)

// creates parsers which match empty test
type EmptyFactory struct{}

// creates directives
type EmptyParser struct{}

// creates parsers which match, in order, the passed strings.
type MatchFactory struct {
	matches []string
}

// creates directives
type MatchParser struct {
	match string
	err   error
	res   postfix.Expression
}

// AnyFactory creates parsers which match any series of lowercase letters
type AnyFactory struct{}

// creates directives
type AnyParser struct{ runes Runes }

// NewExpressionState
func (EmptyFactory) NewExpressionState() ExpressionState           { return EmptyParser{} }
func (EmptyParser) NewRune(rune) State                             { return nil }
func (EmptyParser) GetExpression() (x postfix.Expression, e error) { return }

// NewExpressionState
func (f *MatchFactory) NewExpressionState() ExpressionState {
	var next string
	if cnt := len(f.matches); cnt > 0 {
		next, f.matches = f.matches[0], f.matches[1:]
	}
	return &MatchParser{match: next}
}

func (p MatchParser) GetExpression() (postfix.Expression, error) {
	return p.res, p.err
}

func (p *MatchParser) NewRune(r rune) (ret State) {
	if string(r) == p.match {
		p.res = append(p.res, Quote(p.match))
		ret = Statement(func(rune) State { return nil })
	}
	return
}

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
