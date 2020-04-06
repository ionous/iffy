package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// GetVariable reads from the value at name.
type GetVar struct {
	Name string
}

func (*GetVar) Compose() composer.Spec {
	return composer.Spec{
		Name:  "get_var",
		Spec:  "Get var $1",
		Group: "variables",
		Desc:  "Get Variable: Return the value of the named variable.",
	}
}

func (op *GetVar) GetBool(run rt.Runtime) (ret bool, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = e
	} else {
		ret, err = GetBool(run, p)
	}
	return
}

func (op *GetVar) GetNumber(run rt.Runtime) (ret float64, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = e
	} else {
		ret, err = GetNumber(run, p)
	}
	return
}

func (op *GetVar) GetText(run rt.Runtime) (ret string, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = e
	} else {
		ret, err = GetText(run, p)
	}
	return
}

func (op *GetVar) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = e
	} else {
		ret, err = GetNumbers(run, p)
	}
	return
}

func (op *GetVar) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = e
	} else {
		ret, err = GetTexts(run, p)
	}
	return
}
