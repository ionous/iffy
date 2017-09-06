package express

import (
	"github.com/ionous/iffy/dl/core"
	"go/parser"
	r "reflect"
)

// move to a different package? a sub-package?
// check the imports i guess.
type Xform struct {
	core.Xform
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
	// -- or i guess a command
	if len(t) > 1 {
		panic("hint text eval")
	} else {
		if a, e := parser.ParseExpr(t[0].Str); e != nil {
			err = e
		} else if a, e := ConvertExpr(a, hint); e != nil {
			err = e
		} else {
			ret = a
		}
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
