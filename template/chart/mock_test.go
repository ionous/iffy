package chart

import (
	"github.com/ionous/iffy/template/postfix"
)

// creates parsers
type MockFactory struct{}

// creates directives
type MockParser struct{}

// NewExpressionState
func (MockFactory) NewExpressionState() ExpressionState           { return MockParser{} }
func (MockParser) NewRune(rune) State                             { return nil }
func (MockParser) GetExpression() (x postfix.Expression, e error) { return }

func newText(t string) *TextBlock {
	return &TextBlock{t}
}

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
