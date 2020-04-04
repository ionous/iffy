package next

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type SetVar struct {
	Name rt.TextEval
}

type SetVarBool struct {
	SetVar
	Val rt.BoolEval
}

type SetVarNum struct {
	SetVar
	Val rt.NumberEval
}

type SetVarText struct {
	SetVar
	Val rt.TextEval
}

func (op *SetVar) SetValue(run rt.Runtime, v interface{}) (err error) {
	if name, e := rt.GetText(run, op.Name); e != nil {
		err = e
	} else {
		err = run.SetVariable(name, v)
	}
	return
}

func (*SetVarBool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_bool_var",
		Group: "variables",
		Desc:  "Set Boolean Variable: Sets the named variable to the passed boolean value.",
	}
}

func (op *SetVarBool) Execute(run rt.Runtime) (err error) {
	if val, e := rt.GetBool(run, op.Val); e != nil {
		err = e
	} else {
		err = op.SetValue(run, val)
	}
	return
}

func (*SetVarNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_num_var",
		Group: "variables",
		Desc:  "Set Number Variable: Sets the named variable to the passed number.",
	}
}

func (op *SetVarNum) Execute(run rt.Runtime) (err error) {
	if val, e := rt.GetNumber(run, op.Val); e != nil {
		err = e
	} else {
		err = op.SetValue(run, val)
	}
	return
}

func (*SetVarText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_text_var",
		Group: "variables",
		Desc:  "Set Text Variable: Sets the named variable to the passed piece of text.",
	}
}

func (op *SetVarText) Execute(run rt.Runtime) (err error) {
	if val, e := rt.GetText(run, op.Val); e != nil {
		err = e
	} else {
		err = op.SetValue(run, val)
	}
	return
}
