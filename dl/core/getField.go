package core

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// Field a property value from an object by name.
type Field struct {
	Obj   rt.ObjectEval
	Field string
}

func (*Field) Compose() composer.Spec {
	return composer.Spec{
		Name: "get_field",
		// fix: should use determiner; that should be a hint... yeah?
		// but if so... then doesnt field need that determiner?
		Spec:  "the {field:text_eval} of {object:object_eval}",
		Group: "objects",
		Desc:  "Get Field: Return the value of the named object property.",
	}
}

func (op *Field) GetEval() interface{} {
	return op
}

// GetAssignedValue implements Assignment so we can SetXXX values from variables without a FromXXX statement in between.
func (op *Field) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.getField(run, "")
}

func (op *Field) GetBool(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.Bool)
}

func (op *Field) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.Number)
}

func (op *Field) GetText(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.Text)
}

func (op *Field) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.Text)
}

func (op *Field) GetObject(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.Object)
}

func (op *Field) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.NumList)
}

func (op *Field) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.TextList)
}

func (op *Field) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.RecordList)
}

func (op *Field) getField(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := safe.Field(run, op.Obj, op.Field, aff); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
