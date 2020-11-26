package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

type Pack struct {
	Record rt.RecordEval
	Field  string
	From   Assignment
}

func (*Pack) Compose() composer.Spec {
	return composer.Spec{
		Name:  "pack",
		Group: "variables",
		Desc:  "Pack: Puts a value into a record.",
	}
}

func (op *Pack) Execute(run rt.Runtime) (err error) {
	if obj, e := safe.GetRecord(run, op.Record); e != nil {
		err = e
	} else if val, e := GetAssignedValue(run, op.From); e != nil {
		err = e
	} else {
		// if its going to a record, it should have been a move or copy assignment.
		// in either case, we're overwriting the value.
		err = obj.SetFieldByName(op.Field, val)
	}
	return
}
