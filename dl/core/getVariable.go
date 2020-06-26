package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// GetVariable reads a value of the specified name from the current scope.
// ( ex. loop locals )
type GetVar struct {
	Name string
}

// Compose implements composer.Slat
func (*GetVar) Compose() composer.Spec {
	return composer.Spec{
		Name:  "get_var",
		Spec:  "the {name:text|quote}",
		Group: "variables",
		Desc:  "Get Variable: Return the value of the named variable.",
	}
}

func (op *GetVar) GetBool(run rt.Runtime) (ret bool, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = CmdError{op, e}
	} else if v, e := GetBool(run, p); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetNumber(run rt.Runtime) (ret float64, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = CmdError{op, e}
	} else if v, e := GetNumber(run, p); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetText(run rt.Runtime) (ret string, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = CmdError{op, e}
	} else if v, e := GetText(run, p); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = CmdError{op, e}
	} else if v, e := GetNumbers(run, p); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if p, e := run.GetVariable(op.Name); e != nil {
		err = CmdError{op, e}
	} else if v, e := GetTexts(run, p); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}
