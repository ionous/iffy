package template

import (
	"github.com/ionous/iffy/spec"
)

type NewName interface {
	NewName(string) string
}
type DirectiveParser func(spec.Block, []string) error

type Factory struct {
	gen            NewName
	parseDirective DirectiveParser
}

func MakeFactory(gen NewName, parseDirective DirectiveParser) Factory {
	return Factory{gen, parseDirective}
}

// Tokenize turns a string into a template.
func (f Factory) Tokenize(s string) (ret Template, okay bool) {
	if ts := Tokenize(s); len(ts) > 0 {
		okay, ret = true, Template{ts, f.gen, f.parseDirective}
	}
	return
}
