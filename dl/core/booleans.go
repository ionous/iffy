package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
)

// Is transparently returns a boolean eval.
// It exists to help smooth the use of command expressions:
// eg. c.Cmd("is", "{some expression}")
type Is struct {
	rt.BoolEval
}

// IsNot returns the opposite of a boolean eval.
type IsNot struct {
	rt.BoolEval
}

// IsClass returns true when the object is compatible with the named class.
type IsClass struct {
	Obj   rt.ObjectEval
	Class string
}

// IsExactClass returns true when the object is exactly the named class.
type IsExactClass struct {
	Obj   rt.ObjectEval
	Class string
}

func (op *IsNot) GetBool(run rt.Runtime) (ret bool, err error) {
	if val, e := op.BoolEval.GetBool(run); e != nil {
		err = errutil.New("IsNot.Negate", e)
	} else {
		ret = !val
	}
	return
}

func (op *IsExactClass) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := op.Obj.GetObject(run); e != nil {
		err = e
	} else if cls := obj.Type(); class.IsSame(cls, op.Class) {
		ret = true
	}
	return
}

func (op *IsClass) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := op.Obj.GetObject(run); e != nil {
		err = e
	} else if cls := obj.Type(); class.IsCompatible(cls, op.Class) {
		ret = true
	}
	return
}
