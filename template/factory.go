package template

import (
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/spec"
	r "reflect"
	"strings"
)

type NewName interface {
	NewName(string) string
}
type DirectiveParser func(spec.Block, []string, r.Type) error

type Factory struct {
	gen            NewName
	parseDirective DirectiveParser
}

func MakeFactory(gen NewName, parseDirective DirectiveParser) Factory {
	return Factory{gen, parseDirective}
}

// Templatize turns a string into a template.
func (f Factory) Templatize(c spec.Block, s string) error {
	return f.TemplatizeTokens(c, Tokenize(s))
}

// Templatize turns a string into a template.
func (f Factory) TemplatizeTokens(c spec.Block, ts []Token) (err error) {
	if cnt := len(ts); cnt == 1 {
		ds := strings.Fields(ts[0].Str)
		err = f.parseDirective(c, ds, kindOf.TypeTextEval)
	} else if cnt > 0 {
		ctx := tcontext{f.gen, f.parseDirective}
		err = ctx.convertMulti(c, ts)
	}
	return
}
