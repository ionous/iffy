package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type At struct {
	List  string // variable name
	Index rt.NumberEval
}

// future: lists of lists? probably through lists of records containing lists.
func (*At) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_at",
		Group: "list",
		Spec:  "{list:text} entry {index:number}",
		Desc:  "Value of List: Get a value from a list. The first element is is index 1.",
	}
}

func (op *At) GetNumber(run rt.Runtime) (ret float64, err error) {
	if el, e := op.getEl(run); e != nil {
		err = cmdError(op, e)
	} else if n, e := el.GetNumber(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = float64(n)
	}
	return
}

func (op *At) GetText(run rt.Runtime) (ret string, err error) {
	if el, e := op.getEl(run); e != nil {
		err = cmdError(op, e)
	} else if n, e := el.GetText(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = n
	}
	return
}

func (op *At) GetObject(run rt.Runtime) (ret g.Value, err error) {
	if el, e := op.getEl(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = el
	}
	return
}

//
func (op *At) getEl(run rt.Runtime) (ret g.Value, err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = e
	} else if idx, e := rt.GetNumber(run, op.Index); e != nil {
		err = e
	} else if el, e := vs.GetIndex(int(idx) - 1); e != nil {
		err = e
	} else {
		ret = el
	}
	return
}
