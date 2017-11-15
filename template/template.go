package template

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
)

type Expression = postfix.Expression

// Parse the passed template string into an expression.
func Parse(template string) (ret Expression, err error) {
	p := chart.MakeTemplateParser()
	e := chart.Parse(&p, template)
	xs, ex := p.GetExpression()
	if ex != nil {
		err = errutil.New(ex, e)
	} else if e != nil {
		err = e
	} else {
		ret = xs
	}
	return
}

// ParseExpression reads a series of simple operand and operator phrases.
func ParseExpression(str string) (ret Expression, err error) {
	var p chart.SeriesParser
	if e := chart.Parse(&p, str); e != nil {
		err = e
	} else {
		ret, err = p.GetExpression()
	}
	return
}
