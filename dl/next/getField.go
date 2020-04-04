package next

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// GetField a property value from an object by name.
type GetField struct {
	Obj, Field rt.TextEval
}

func (*GetField) Compose() composer.Spec {
	return composer.Spec{
		Name:  "get_field",
		Spec:  "Get $2 of $1",
		Group: "objects",
		Desc:  "Get Field: Return the value of the named object property.",
	}
}

func (op *GetField) GetBool(run rt.Runtime) (ret bool, err error) {
	err = op.GetValue(run, &ret)
	return
}

func (op *GetField) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = op.GetValue(run, &ret)
	return
}

func (op *GetField) GetText(run rt.Runtime) (ret string, err error) {
	err = op.GetValue(run, &ret)
	return
}

func (op *GetField) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	var values []float64
	if e := op.GetValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewNumberList(values)
	}
	return
}

func (op *GetField) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	var values []string
	if e := op.GetValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewTextList(values)
	}
	return
}

func (op *GetField) GetValue(run rt.Runtime, pv interface{}) (err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if field, e := rt.GetText(run, op.Field); e != nil {
		err = e
	} else {
		err = run.GetField(obj, field, pv)
	}
	return
}
