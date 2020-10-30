package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

// GetVariable reads a value of the specified name from the current scope.
// ( ex. loop locals, or -- in a noun scope -- might translate "apple" to "$macintosh" )
type GetVar struct {
	Name  rt.TextEval // uses text eval to make template expressions easier
	Flags TryAsNoun   `if:"internal"`
}

type UnknownVariable string

func (e UnknownVariable) Error() string {
	return errutil.Sprintf("Unknown variable %q", string(e))
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
	if local, e := getVariableByName(run, op.Name, op.Flags); e != nil {
		err = cmdError(op, e)
	} else if v, e := local.GetBool(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetNumber(run rt.Runtime) (ret float64, err error) {
	if local, e := getVariableByName(run, op.Name, op.Flags); e != nil {
		err = cmdError(op, e)
	} else if v, e := local.GetNumber(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetText(run rt.Runtime) (ret string, err error) {
	if local, e := getVariableByName(run, op.Name, op.Flags); e != nil {
		err = cmdError(op, e)
	} else if vs, e := local.GetText(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs
	}
	return

}

// allows us to use GetVar directly in things that take an object.
func (op *GetVar) GetObjectValue(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getObjectValue(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetVar) getObjectValue(run rt.Runtime) (ret rt.Value, err error) {
	local, e := getVariableByName(run, op.Name, op.Flags)
	switch e := e.(type) {
	case nil:
		ret = local
	default:
		err = e
	case UnknownVariable:
		if op.Flags.tryObject() {
			ret, err = getObjectExactly(run, string(e))
		}
	}
	return
}

func (op *GetVar) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if local, e := getVariableByName(run, op.Name, op.Flags); e != nil {
		err = cmdError(op, e)
	} else if vs, e := local.GetNumList(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs
	}
	return
}

func (op *GetVar) GetTextList(run rt.Runtime) (ret []string, err error) {
	if local, e := getVariableByName(run, op.Name, op.Flags); e != nil {
		err = cmdError(op, e)
	} else if vs, e := local.GetTextList(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs
	}
	return
}

// GetVar asks for a variable using a text eval;
// we first need to determine which actual variable name they mean.
func getVariableByName(run rt.Runtime, text rt.TextEval, flags TryAsNoun) (ret rt.Value, err error) {
	// first resolve the requested variable name into text
	if n, e := rt.GetText(run, text); e != nil {
		err = e
	} else {
		if !flags.tryVariable() {
			err = UnknownVariable(n)
		} else {
			switch v, e := run.GetField(object.Variables, n); e.(type) {
			case nil:
				ret = v
			default:
				err = e
			case rt.UnknownTarget, rt.UnknownField:
				err = UnknownVariable(n)
			}
		}
	}
	return
}
