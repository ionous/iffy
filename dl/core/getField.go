package core

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// GetField a property value from an object by name.
type GetField struct {
	Obj   rt.ObjectEval
	Field string
}

func (*GetField) Compose() composer.Spec {
	return composer.Spec{
		Name: "get_field",
		// fix: should use determiner; that should be a hint... yeah?
		// but if so... then doesnt GetField need that determiner?
		Spec:  "the {field:text_eval} of {object:object_ref}",
		Group: "objects",
		Desc:  "Get Field: Return the value of the named object property.",
	}
}

func (op *GetField) GetEval() interface{} {
	return op
}

// GetAssignedValue implements Assignment so we can SetXXX values from variables without a FromXXX statement in between.
func (op *GetField) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.getField(run, "")
}

func (op *GetField) GetBool(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.Bool)
}

func (op *GetField) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.Number)
}

func (op *GetField) GetText(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.Text)
}

func (op *GetField) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.NumList)
}

func (op *GetField) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.getField(run, affine.TextList)
}

func (op *GetField) getField(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := safe.Field(run, op.Obj, op.Field, aff); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
