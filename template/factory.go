package template

import (
	"github.com/ionous/iffy/spec"
)

type GenerateId interface {
	GenerateId(string) string
}
type ExpressionParser func(spec.Block, string) error

type Factory struct {
	gen       GenerateId
	parseExpr ExpressionParser
}

func MakeFactory(gen GenerateId, parseExpr ExpressionParser) Factory {
	return Factory{gen, parseExpr}
}

// Tokenize turns a string into a template.
func (f Factory) Tokenize(s string) (ret Template, okay bool) {
	if ts, cnt := tokenize(s); cnt > 0 {
		okay, ret = true, Template{ts, f.gen, f.parseExpr}
	}
	return
}
