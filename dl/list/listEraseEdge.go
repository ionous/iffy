package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

/*
 * erase: numEval
 * from: varName,
 * atIndex: num,
 */
type EraseEdge struct {
	From   ListSource `if:"unlabeled"`
	AtEdge Edge       `if:"unlabeled"`
}

func (*EraseEdge) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erase", Role: composer.Command},
		Desc:   "Erase: Remove one or more values from a list",
	}
}

func (op *EraseEdge) Execute(run rt.Runtime) (err error) {
	if e := op.pop(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *EraseEdge) pop(run rt.Runtime) (err error) {
	if vs, e := GetListSource(run, op.From); e != nil {
		err = e
	} else {
		if cnt := vs.Len(); cnt > 0 {
			var at int
			if !op.AtEdge.Front() {
				at = cnt - 1
			}
			if _, e := vs.Splice(at, at+1, nil); e != nil {
				err = e
			}
		}
	}
	return
}
