package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Push struct {
	List   string // variable name
	Insert core.Assignment
	Front  *Front
}

func (op *Push) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_push",
		Group: "list",
		Spec:  "push {into%list:text} {front?list_edge} {inserting%insert:assignment}",
		Desc: `Push into list: Add elements to the front or back of a list.
Returns the new length of the list.`,
	}
}

func (op *Push) Execute(run rt.Runtime) (err error) {
	if _, e := op.push(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

// returns the new size
func (op *Push) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if cnt, e := op.push(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.IntOf(cnt)
	}
	return
}

func (op *Push) push(run rt.Runtime) (ret int, err error) {
	if els, e := safe.List(run, op.List); e != nil {
		err = e
	} else if ins, e := core.GetAssignedValue(run, op.Insert); e != nil {
		err = e
	} else if !IsAppendable(ins, els) {
		err = insertError{ins, els}
	} else {
		if op.Front == nil || !*op.Front {
			els.Append(ins)
		} else {
			_, err = els.Splice(0, 0, ins)
		}
		if err == nil {
			ret = els.Len()
		}
	}
	return
}
