package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	r "reflect"
)

type exform struct {
	cmds      *ops.Ops
	templates template.Factory
}

func MakeXform(cmds *ops.Ops, gen template.NewName) ops.Transform {
	return exform{cmds, template.MakeFactory(gen, ParseDirective)}
}

// TransformValue returns src if no error but couldnt convert.
func (xf exform) TransformValue(v r.Value, hint r.Type) (ret r.Value, err error) {
	base := core.Xform{}
	if t, ok := tryTokenize(xf.templates, v); !ok {
		ret, err = base.TransformValue(v, hint)
	} else {
		target := ops.NewValue(hint)
		c := xf.cmds.NewFromTarget(target, base)
		err = t.Convert(c)
		if err == nil {
			if e := c.Build(); e != nil {
				err = e
			} else {
				ret = target.Field(0)
			}
		}
	}
	return
}

// tryTokenize attempt to turn the passed val into a string template.
func tryTokenize(f template.Factory, val r.Value) (ret template.Template, okay bool) {
	if val.Kind() == r.String {
		ret, okay = f.Tokenize(val.String())
	}
	return
}
