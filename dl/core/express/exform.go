package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
	"go/parser"
	r "reflect"
)

type Xform struct {
	cmds *ops.Ops
	core.Xform
}

func MakeXform(cmds *ops.Ops) Xform {
	return Xform{cmds: cmds}
}

// TransformValue returns src if no error but couldnt convert.
func (ts Xform) TransformValue(val interface{}, hint r.Type) (ret interface{}, err error) {
	if t, ok := tryTokenize(val); ok {
		ret, err = ts.TransformTemplate(t, hint)
	} else {
		ret, err = ts.Xform.TransformValue(val, hint)
	}
	return
}

// TransformTemplate
func (ts Xform) TransformTemplate(t Template, hint r.Type) (ret interface{}, err error) {
	// FIX: not just one token? than our output sure better be a text eval
	if len(t) > 1 {
		for _, x := range t {
			println(x.Str)
		}
		panic("hint text eval")
	} else {
		// look for and chomp templates starting with {go}
		if g, ok := t[0].Go(); ok {
			// create a factory for new building comands
			fac := ops.NewFactory(ts.cmds, ts)
			// use the raw interface for building commands
			// b/c we know that we only want to build commands
			// but the formatting of them are linear in a single string
			// rather than spread out in iffy blocks.
			if cmd, e := fac.NewSpec(g[0]); e != nil {
				err = e
			} else if e := converts(cmd, g[1:]); e != nil {
				err = e
			} else {
				// get the underlying value that ops created.
				ret = cmd.(*ops.Command).Target().Interface()
			}
		} else if a, e := convert(t[0].Str); e != nil {
			err = e
		} else {
			ret = a
		}
	}
	return
}

// add the passed strings as parsed expressions to the passed cmd
// we could probably support keywords using name:something in the string.
func converts(cmd spec.Spec, strs []string) (err error) {
	for _, s := range strs {
		if a, e := convert(s); e != nil {
			err = e
			break
		} else if e := cmd.Position(a); e != nil {
			err = e
			break
		}
	}
	return
}

func convert(s string) (ret interface{}, err error) {
	if a, e := parser.ParseExpr(s); e != nil {
		err = e
	} else if a, e := ConvertExpr(a); e != nil {
		err = e
	} else {
		ret = a
	}
	return
}

// tryTokenize attempt to turn the passed val into a string template.
func tryTokenize(val interface{}) (ret Template, okay bool) {
	if s, ok := val.(string); ok {
		ret, okay = Tokenize(s)
	}
	return
}
