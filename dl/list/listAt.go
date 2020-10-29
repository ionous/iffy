package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

type At struct {
	List  string // variable name
	Index int
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
	if vs, e := op.getValue(run); e != nil {
		err = cmdError(op, e)
	} else if n, e := vs.GetNumber(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = float64(n)
	}
	return
}

func (op *At) GetText(run rt.Runtime) (ret string, err error) {
	if vs, e := op.getValue(run); e != nil {
		err = cmdError(op, e)
	} else if n, e := vs.GetText(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = n
	}
	return
}

func (op *At) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if vs, e := op.getValue(run); e != nil {
		err = cmdError(op, e)
	} else if vs, e := vs.GetNumList(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs
	}
	return
}

func (op *At) GetTextList(run rt.Runtime) (ret []string, err error) {
	if vs, e := op.getValue(run); e != nil {
		err = cmdError(op, e)
	} else if vs, e := vs.GetTextList(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs
	}
	return
}

func (op *At) getValue(run rt.Runtime) (ret rt.Value, err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = e
	} else if el, e := vs.GetIndex(op.Index); e != nil {
		err = e
	} else {
		ret = el
	}
	return
}
