package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type SetField struct {
	Obj        ObjectRef
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
	if id, e := GetObjectRef(run, op.Obj); e != nil {
		err = e
	} else if field, e := rt.GetText(run, op.Field); e != nil {
		err = e
	} else if val, e := GetAssignedValue(run, op.Assignment); e != nil {
		err = e
	} else {
		err = run.SetField(id, field, val)
	}
	return
}
