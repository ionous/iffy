package render

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type RenderTemplate struct {
	Expression rt.TextEval `if:"internal"`
}

func (*RenderTemplate) Compose() composer.Spec {
	return composer.Spec{
		Spec:  "the template {lines%template:lines|quote}",
		Group: "format",
		Desc:  "Render Template: Parse text using iffy templates. See: https://github.com/ionous/iffy/wiki/Templates",
		Stub:  true,
	}
}

// RunTest returns an error on failure.
func (op *RenderTemplate) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetOptionalText(run, op.Expression, ""); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
