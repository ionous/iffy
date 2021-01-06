package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

/**
 * put: eval(num,txt,rec),
 * intoNum/Txt/RecList: varName,
 * atBack|atFront.
 */
type PutEdge struct {
	From   core.Assignment `if:"selector"`
	Into   ListTarget      `if:"selector"`
	AtEdge Edge            `if:"selector"`
}

func (*PutEdge) Compose() composer.Spec {
	return composer.Spec{
		Name:   "put_edge",
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Desc:   "Put: add a value to a list",
	}
}

func (op *PutEdge) Execute(run rt.Runtime) (err error) {
	if e := op.push(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *PutEdge) push(run rt.Runtime) (err error) {
	if ins, e := core.GetAssignedValue(run, op.From); e != nil {
		err = e
	} else if els, e := op.Into.GetListTarget(run); e != nil {
		err = e
	} else if !IsAppendable(ins, els) {
		err = insertError{ins, els}
	} else {
		if !op.AtEdge.Front() {
			els.Append(ins)
		} else {
			_, err = els.Splice(0, 0, ins)
		}
	}
	return
}
