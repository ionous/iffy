package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	r "reflect"
)

type exform struct {
	cmds *ops.Ops
	facs template.Factory
}

func MakeXform(cmds *ops.Ops, gen template.NewName) ops.Transform {
	facs := template.MakeFactory(gen, ParseDirective)
	return exform{cmds, facs}
}

// TransformValue returns src if no error but couldnt convert.
func (xf exform) TransformValue(v r.Value, hint r.Type) (ret r.Value, err error) {
	base := core.Xform{}
	if ts, ok := tryTokenize(v); !ok {
		ret, err = base.TransformValue(v, hint)
	} else {
		target := ops.NewValue(hint)
		b := xf.cmds.NewFromTarget(target, base)
		// c := clog.Make(os.Stderr, b)
		c := b
		if e := xf.facs.TemplatizeTokens(c, ts); e != nil {
			err = e
		} else if e := b.Build(); e != nil {
			err = e
		} else {
			ret = target.Field(0)
		}
	}
	return
}

// tryTokenize attempt to turn the passed val into a string template.
func tryTokenize(val r.Value) (ret []template.Token, okay bool) {
	if val.Kind() == r.String {
		ret, okay = template.Tokenize(val.String()), true
	}
	return
}
