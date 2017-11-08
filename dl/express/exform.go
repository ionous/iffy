package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/text"
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
	if ds, e := tryTokenize(v); e != nil {
		err = e
	} else if cnt := len(ds); cnt == 0 {
		// no directives? transform the value via core.
		ret, err = core.Transform(v, hint)
	} else {
		dir := ds[0]
		cmdf := ops.NewFactory(xf.cmds, nil)
		if cnt == 1 && len(dir.Key) == 0 {
			// one directive; transform the expression directly.
			if cmd, e := Convert(cmdf, dir.Expression); e != nil {
				err = e
			} else {
				ret = cmd.Target()
			}
		} else {
			// multiple directives? then we need more firepower:
			// convert via "text" which manages the structure of
			em := expressionMaker{cmdf, xf.gen}
			if cmd, e := text.ConvertDirectives(em, ds); e != nil {
				err = e
			} else {
				ret = cmd.Target()
			}
		}
	}
	return
}

// tryTokenize attempt to turn the passed val into a string template.
func tryTokenize(val r.Value) (ret []template.Directive, err error) {
	if val.Kind() == r.String {
		if str := val.String(); strings.Contains(str, "{") {
			ret, err = chart.Parse(str)
		}
	}
	return
}

type expressionMaker struct {
	*ops.Factory
	gen ident.Counters
}

func (tm expressionMaker) CreateName(group string) (string, error) {
	return tm.gen.NewName(group), nil
}

func (tm expressionMaker) CreateExpression(xs postfix.Expression, hint r.Type) (ret *ops.Command, err error) {
	// FIX: hint?
	return Convert(tm.Factory, xs)
}
