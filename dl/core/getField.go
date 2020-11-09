package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// GetField a property value from an object by name.
type GetField struct {
	Obj   ObjectEval
	Field rt.TextEval
}

func (*GetField) Compose() composer.Spec {
	return composer.Spec{
		Name: "get_field",
		// fix: should use determiner; that should be a hint... yeah?
		// but if so... then doesnt GetField need that determiner?
		Spec:  "the {field:text_eval} of {object%obj:object_ref}",
		Group: "objects",
		Desc:  "Get Field: Return the value of the named object property.",
	}
}

func (op *GetField) GetBool(run rt.Runtime) (ret bool, err error) {
	if p, e := op.getValue(run); e != nil {
		err = e
	} else {
		ret, err = p.GetBool()
	}
	return
}

func (op *GetField) GetNumber(run rt.Runtime) (ret float64, err error) {
	if p, e := op.getValue(run); e != nil {
		err = e
	} else {
		ret, err = p.GetNumber()
	}
	return
}

func (op *GetField) GetText(run rt.Runtime) (ret string, err error) {
	if p, e := op.getValue(run); e != nil {
		err = e
	} else {
		ret, err = p.GetText()
	}
	return
}

func (op *GetField) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if p, e := op.getValue(run); e != nil {
		err = e
	} else {
		ret, err = p.GetNumList()
	}
	return
}

func (op *GetField) GetTextList(run rt.Runtime) (ret []string, err error) {
	if p, e := op.getValue(run); e != nil {
		err = e
	} else {
		ret, err = p.GetTextList()
	}
	return
}

func (op *GetField) getValue(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := GetObjectValue(run, op.Obj); e != nil {
		err = e
	} else if field, e := rt.GetText(run, op.Field); e != nil {
		err = e
	} else {
		ret, err = obj.GetNamedField(field)
	}
	return
}
