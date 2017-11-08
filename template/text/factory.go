package text

import (
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
	r "reflect"
)

type Factory interface {
	CreateName(string) (string, error)
	CreateExpression(spec.Block, postfix.Expression, r.Type) error
}

func String(f Factory, c spec.Block, s string) (err error) {
	if dirs, e := chart.Parse(s); e != nil {
		err = e
	} else if cnt := len(dirs); cnt > 0 {
		err = Directives(f, c, dirs)
	}
	return
}

func Directives(f Factory, c spec.Block, dirs []chart.Directive) error {
	return Engine{f}.reduce(c, dirs)
}
