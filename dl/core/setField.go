package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type SetField struct {
	Obj   rt.ObjectEval
	Field string
	From  Assignment
}

func (*SetField) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_field",
		Group: "objects",
		Desc:  "Set Field: Sets the named field to the assigned value.",
	}
}

func (op *SetField) Execute(run rt.Runtime) (err error) {
	if obj, e := rt.GetObject(run, op.Obj); e != nil {
		err = e
	} else if val, e := GetAssignedValue(run, op.From); e != nil {
		err = e
	} else {
		// if its going to a record, it should have been a move or copy assignment.
		// in either case, we're overwriting the value.
		err = obj.SetNamedField(op.Field, val)
	}
	return
}
