package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
)

type IsNot struct {
	Negate rt.BoolEval
}

type IsSameClass struct {
	Obj   rt.ObjectEval
	Class string
}

type IsSimilarClass struct {
	Obj   rt.ObjectEval
	Class string
}

func (op *IsNot) GetBool(run rt.Runtime) (ret bool, err error) {
	if val, e := op.Negate.GetBool(run); e != nil {
		err = errutil.New("IsNot.Negate", e)
	} else {
		ret = !val
	}
	return
}

func (op *IsSameClass) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := op.Obj.GetObject(run); e != nil {
		err = e
	} else if cls, ok := run.GetClass(op.Class); ok {
		ret = cls == obj.Type()
	}
	return
}

func (op *IsSimilarClass) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := op.Obj.GetObject(run); e != nil {
		err = e
	} else if cls := obj.Type(); class.IsCompatible(cls, op.Class) {
		ret = true
	}
	return
}
