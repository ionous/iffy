package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

type Set struct {
	List    string // variable name
	Index   rt.NumberEval
	Element core.Assignment
}

func (*Set) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_set",
		Group: "list",
		Desc:  "Set Value of List: Overwrite an existing value in a list.",
	}
}

func (op *Set) Execute(run rt.Runtime) (err error) {
	if e := op.setAt(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Set) setAt(run rt.Runtime) (err error) {
	if els, e := run.GetField(object.Variables, op.List); e != nil {
		err = e
	} else if onedex, e := rt.GetNumber(run, op.Index); e != nil {
		err = e
	} else if el, e := core.GetAssignedValue(run, op.Element); e != nil {
		err = e
	} else {
		i := int(onedex - 1)
		err = els.SetIndexedValue(i, el)
	}
	return
}
