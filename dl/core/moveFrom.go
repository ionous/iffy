package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type MoveFrom struct {
	Name  string    // name of the variable or object.
	Flags TryAsNoun `if:"internal"`
}

func (*MoveFrom) Compose() composer.Spec {
	return composer.Spec{
		Name:  "move_from",
		Group: "variables",
		Desc:  `Move Variable: Move the contents of one variable to another, leaving the first variable blank.`,
	}
}

func (op *MoveFrom) GetEval() interface{} { return nil }
func (op *MoveFrom) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := op.moveFrom(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (op *MoveFrom) moveFrom(run rt.Runtime) (ret g.Value, err error) {
	if box, val, e := getVariableValue(run, op.Name, op.Flags); e != nil {
		err = e
	} else {
		if val != nil {
			// we found the named variable;
			// clear out the old contents
			ret, err = val, box.SetValue(run, nil)
		} else if op.Flags.tryObject() {
			// its an object reference, move is the same as copy.
			ret, err = box.GetObjectByName(run)
		} else {
			err = g.UnknownObject(string(box))
		}
	}
	return
}
