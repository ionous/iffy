package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	r "reflect"
	"strings"
)

type exform struct {
	cmds *ops.Ops
	gen  ident.Counters
}

// NewTransform returns a transform which can convert command text into actions.
// For example, in c.Cmd("set text", "status", "{score}/{turnCount}") --
// the "{}" would be converted into a tree of iffy commands.
func NewTransform(cmds *ops.Ops, gen ident.Counters) *exform {
	return &exform{cmds, gen}
}

// TransformValue returns src if no error but couldnt convert.
func (xf *exform) TransformValue(v r.Value, hint r.Type) (ret r.Value, err error) {
	if xs, e := tryTokenize(v); e != nil {
		err = e
	} else if cnt := len(xs); cnt == 0 {
		// no directives? transform the value via core ( for literals )
		ret, err = core.Transform(v, hint)
	} else {
		cmds := ops.NewFactory(xf.cmds, nil)
		if cmd, e := Convert(cmds, xs, xf.gen); e != nil {
			err = e
		} else {
			ret = cmd.Target()
		}
	}
	return
}

// tryTokenize attempt to turn the passed val into a string types.
func tryTokenize(val r.Value) (ret template.Expression, err error) {
	if val.Kind() == r.String {
		if str := val.String(); strings.Contains(str, "{") {
			ret, err = template.Parse(str)
		}
	}
	return
}
