package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type SetField struct {
	Obj        ObjectEval
	Field      rt.TextEval
	Assignment Assignment
}

func (*SetField) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_field",
		Group: "objects",
		Desc:  "Set Field: Sets the named field to the assigned value.",
	}
}

func (op *SetField) Execute(run rt.Runtime) (err error) {
	if obj, e := GetObjectValue(run, op.Obj); e != nil {
		err = e
	} else if field, e := rt.GetText(run, op.Field); e != nil {
		err = e
	} else if val, e := GetAssignedValue(run, op.Assignment); e != nil {
		err = e
	} else {
		// if its going to a record, it should have been a move or copy assignment.
		// in either case, we're overwriting the value.
		err = obj.SetNamedField(field, val)
	}
	return
}
