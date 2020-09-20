package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// GetVariable reads a value of the specified name from the current scope.
// ( ex. loop locals, or -- in a noun scope -- might translate "apple" to "$macintosh" )
type GetVar struct {
	Name rt.TextEval // uses text eval to make template expressions easier
}

// Compose implements composer.Slat
func (*GetVar) Compose() composer.Spec {
	return composer.Spec{
		Name:  "get_var",
		Spec:  "the {name:text_eval}",
		Group: "variables",
		Desc:  "Get Variable: Return the value of the named variable.",
	}
}

func (op *GetVar) GetBool(run rt.Runtime) (ret bool, err error) {
	if p, e := op.getVar(run); e != nil {
		err = CmdError{op, e}
	} else if v, e := p.GetBool(run); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetNumber(run rt.Runtime) (ret float64, err error) {
	if p, e := op.getVar(run); e != nil {
		err = CmdError{op, e}
	} else if v, e := p.GetNumber(run); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetText(run rt.Runtime) (ret string, err error) {
	if p, e := op.getVar(run); e != nil {
		err = CmdError{op, e}
	} else if v, e := p.GetText(run); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if p, e := op.getVar(run); e != nil {
		err = CmdError{op, e}
	} else if v, e := p.GetNumberStream(run); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if p, e := op.getVar(run); e != nil {
		err = CmdError{op, e}
	} else if v, e := p.GetTextStream(run); e != nil {
		err = CmdError{op, e}
	} else {
		ret = v
	}
	return
}

func (op *GetVar) getVar(run rt.Runtime) (ret rt.Value, err error) {
	if n, e := rt.GetText(run, op.Name); e != nil {
		err = e
	} else {
		ret, err = run.GetVariable(n)
	}
	return
}
