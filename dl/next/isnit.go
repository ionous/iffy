package next

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

// Is transparently returns a boolean eval.
// It exists to help smooth the use of command expressions:
// eg. "is" {some expression}
type Is struct {
	Test rt.BoolEval
}

// IsNot returns the opposite of a boolean eval.
type IsNot struct {
	Test rt.BoolEval
}

func (op *Is) GetBool(run rt.Runtime) (ret bool, err error) {
	if val, e := rt.GetBool(run, op.Test); e != nil {
		err = errutil.New("IsNot.Negate", e)
	} else {
		ret = val
	}
	return
}

func (op *IsNot) GetBool(run rt.Runtime) (ret bool, err error) {
	if val, e := rt.GetBool(run, op.Test); e != nil {
		err = errutil.New("IsNot.Negate", e)
	} else {
		ret = !val
	}
	return
}
