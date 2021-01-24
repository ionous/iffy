package core

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// GetAtField a property value from an object by name.
type GetAtField struct {
	Field string           `if:"selector"`
	From  FromSourceFields `if:"selector"`
}

func (*GetAtField) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "get", Role: composer.Function},
		Group:  "variables",
		Desc:   "GetAtField: Get a value from a record.",
	}
}

func (op *GetAtField) Affinity() affine.Affinity { return "" }

// GetAssignedValue implements Assignment so we can SetXXX values from variables without a FromXXX statement in between.
func (op *GetAtField) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, "")
}

func (op *GetAtField) GetBool(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Bool)
}

func (op *GetAtField) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Number)
}

func (op *GetAtField) GetText(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Text)
}

func (op *GetAtField) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Record)
}

func (op *GetAtField) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.NumList)
}

func (op *GetAtField) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.TextList)
}

func (op *GetAtField) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.RecordList)
}

func (op *GetAtField) unpack(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if src, e := GetSourceFields(run, op.From); e != nil {
		err = cmdError(op, e)
	} else if v, e := safe.Unpack(src, op.Field, aff); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
