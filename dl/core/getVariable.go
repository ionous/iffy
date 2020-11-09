package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// GetVariable reads a value of the specified name from the current scope.
// ( ex. loop locals, or -- in a noun scope -- might translate "apple" to "$macintosh" )
type GetVar struct {
	Name  rt.TextEval // uses text eval to make template expressions easier
	Flags TryAsNoun   `if:"internal"`
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
	if box, e := getVariableNamed(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if res, e := box.GetBool(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = res
	}
	return
}

func (op *GetVar) GetNumber(run rt.Runtime) (ret float64, err error) {
	if box, e := getVariableNamed(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if res, e := box.GetNumber(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = res
	}
	return
}

func (op *GetVar) GetText(run rt.Runtime) (ret string, err error) {
	if box, e := getVariableNamed(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if res, e := box.GetText(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = res
	}
	return

}

// allows us to use GetVar directly in commands which take a named object.
func (op *GetVar) GetObjectValue(run rt.Runtime) (ret g.Value, err error) {
	if box, e := op.getObjectValue(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = box
	}
	return
}

func (op *GetVar) getObjectValue(run rt.Runtime) (ret g.Value, err error) {
	if name, val, e := getVariableValue(run, op.Name, op.Flags); e != nil {
		err = e
	} else if val != nil {
		ret = val
	} else if op.Flags.tryObject() {
		ret, err = name.GetObjectByName(run)
	} else {
		err = rt.UnknownObject(string(name))
	}
	return
}

func (op *GetVar) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if box, e := getVariableNamed(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if res, e := box.GetNumList(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = res
	}
	return
}

func (op *GetVar) GetTextList(run rt.Runtime) (ret []string, err error) {
	if box, e := getVariableNamed(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if res, e := box.GetTextList(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = res
	}
	return
}

func (op *GetVar) GetRecordList(run rt.Runtime) (ret []*g.Record, err error) {
	if box, e := getVariableNamed(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if res, e := box.GetRecordList(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = res
	}
	return
}

// resolve a requested variable name into a helper string which can access its contents
func getVariableNamed(run rt.Runtime, text rt.TextEval) (ret rt.Variable, err error) {
	if n, e := rt.GetText(run, text); e != nil {
		err = e
	} else {
		ret = rt.Variable(n)
	}
	return
}

// return the name, and optionally the value of a named variable
// returns a nil value if the variable couldnt be found
// returns error only critical errors
func getVariableValue(run rt.Runtime, text rt.TextEval, flags TryAsNoun) (retBox rt.Variable, retVal g.Value, err error) {
	// first resolve the requested variable name into text
	if box, e := getVariableNamed(run, text); e != nil {
		err = e
	} else {
		retBox = box
		if flags.tryVariable() {
			switch val, e := box.GetValue(run); e.(type) {
			case nil:
				retVal = val
			default:
				err = e
			case rt.UnknownTarget, rt.UnknownField:
				retVal = nil // no such variable....
			}
		}
	}
	return
}
