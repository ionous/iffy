package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type Make struct {
	Name      string
	Arguments *Arguments // kept as a pointer for composer formatting...
}

func (*Make) Compose() composer.Spec {
	return composer.Spec{
		Name: "make",
	}
}

func (op *Make) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	if d, e := op.makeRecord(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.RecordOf(d)
	}
	return
}

func (op *Make) makeRecord(run rt.Runtime) (ret *g.Record, err error) {
	if k, e := run.GetKindByName(op.Name); e != nil {
		err = e
	} else {
		b := k.NewRecord()
		if op.Arguments == nil {
			ret = b // return the empty record
		} else if e := op.Arguments.Distill(run, nil, b); e != nil {
			err = e
		} else {
			ret = b
		}
	}
	return
}
