package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
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
func (op *GetField) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if p, e := op.getField(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = p
	}
	return
}

func (op *GetField) GetBool(run rt.Runtime) (ret bool, err error) {
	if p, e := op.getField(run); e != nil {
		err = e
	} else {
		ret, err = p.GetBool()
	}
	return
}

func (op *GetField) GetNumber(run rt.Runtime) (ret float64, err error) {
	if p, e := op.getField(run); e != nil {
		err = e
	} else {
		ret, err = p.GetNumber()
	}
	return
}

func (op *GetField) GetText(run rt.Runtime) (ret string, err error) {
	if p, e := op.getField(run); e != nil {
		err = e
	} else {
		ret, err = p.GetText()
	}
	return
}

func (op *GetField) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if p, e := op.getField(run); e != nil {
		err = e
	} else {
		ret, err = p.GetNumList()
	}
	return
}

func (op *GetField) GetTextList(run rt.Runtime) (ret []string, err error) {
	if p, e := op.getField(run); e != nil {
		err = e
	} else {
		ret, err = p.GetTextList()
	}
	return
}

func (op *GetField) getField(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := rt.GetObject(run, op.Obj); e != nil {
		err = e
	} else {
		ret, err = obj.GetNamedField(op.Field)
	}
	return
}
