package chart

import (
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

// Section joins a bunch of separate postfix expressions into a single operation.
// if it only has one expression, thats all thats returned; otherwise a span does.
type Section struct {
	list []postfix.Expression
}

// Add the passed text as an expression to the current span.
func (x *Section) AddText(t string) {
	x.Append(quote(t))
}

// Append the expression to the current span.
func (x *Section) Append(xs postfix.Expression) {
	if len(xs) > 0 {
		x.list = append(x.list, xs)
	}
}

// Reduce returns the summation of the current span.
func (x Section) Reduce(kind template.BuiltinType) (ret postfix.Expression) {
	if cnt := len(x.list); cnt == 1 && kind == template.Span {
		ret = x.list[0]
	} else /*if cnt > 0*/ {
		for _, v := range x.list {
			ret = append(ret, v...)
		}
		ret = append(ret, template.Builtin{kind, len(x.list)})
	}
	return
}
