package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type CopyFrom struct {
	Name  rt.TextEval // name of the variable or object.
	Flags TryAsNoun   `if:"internal"`
}

func (*CopyFrom) Compose() composer.Spec {
	return composer.Spec{
		Name:  "copy_from",
		Group: "variables",
		Desc:  `Copy Variable: Copy the contents of one variable to another.`,
	}
}

func (op *CopyFrom) GetEval() interface{} { return nil }
func (op *CopyFrom) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.copyFrom(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *CopyFrom) copyFrom(run rt.Runtime) (ret g.Value, err error) {
	if box, val, e := getVariableValue(run, op.Name, op.Flags); e != nil {
		err = e
	} else if val != nil {
		ret = val
	} else if op.Flags.tryObject() {
		ret, err = box.GetObjectByName(run)
	} else {
		err = rt.UnknownObject(string(box))
	}
	return
}
