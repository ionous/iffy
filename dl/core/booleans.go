package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type IsNot struct {
	Negate rt.BoolEval
}

func (op *IsNot) GetBool(run rt.Runtime) (ret bool, err error) {
	if val, e := op.Negate.GetBool(run); e != nil {
		err = errutil.New("IsNot.Negate", e)
	} else {
		ret = !val
	}
	return
}
