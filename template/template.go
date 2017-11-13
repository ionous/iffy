package template

import (
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
)

// Parse the passed template string into an expression.
func Parse(template string) (ret postfix.Expression, err error) {
	p := chart.MakeTemplateParser()
	e := chart.Parse(&p, template)
	xs, ex := p.GetExpression()
	if ex != nil {
		err = ex // prefer the expression error if there is one.
	} else if e != nil {
		err = e
	} else {
		ret = xs
	}
	return
}

// ParseExpression reads a series of simple operand and operator phrases.
func ParseExpression(str string) (ret postfix.Expression, err error) {
	var p chart.SeriesParser
	if e := chart.Parse(&p, str); e != nil {
		err = e
	} else {
		ret, err = p.GetExpression()
	}
	return
}
