package optional

import (
	"github.com/ionous/iffy/rt"
)

func Bool(run rt.Runtime, eval rt.BoolEval) (ret bool, err error) {
	if eval != nil {
		if ok, e := eval.GetBool(run); e != nil {
			err = e
		} else {
			ret = ok
		}
	}
	return
}

func Text(run rt.Runtime, text rt.TextEval) (ret string, err error) {
	if text != nil {
		if s, e := text.GetText(run); e != nil {
			err = e
		} else {
			ret = s
		}
	}
	return
}
