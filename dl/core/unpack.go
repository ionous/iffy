package core

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// Unpack a property value from an object by name.
type Unpack struct {
	Record rt.RecordEval
	Field  string
}

func (*Unpack) Compose() composer.Spec {
	return composer.Spec{
		Name:  "unpack",
		Spec:  "unpack {field:text_eval} from {record:record_eval}",
		Group: "variables",
		Desc:  "Unpack: Get a value from a record.",
	}
}

func (op *Unpack) GetEval() interface{} {
	return op
}

// GetAssignedValue implements Assignment so we can SetXXX values from variables without a FromXXX statement in between.
func (op *Unpack) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, "")
}

func (op *Unpack) GetBool(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Bool)
}

func (op *Unpack) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Number)
}

func (op *Unpack) GetText(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Text)
}

func (op *Unpack) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Text)
}

func (op *Unpack) GetObject(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Object)
}

func (op *Unpack) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.NumList)
}

func (op *Unpack) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.TextList)
}

func (op *Unpack) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.RecordList)
}

func (op *Unpack) unpack(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := safe.Unpack(run, op.Record, op.Field, aff); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
