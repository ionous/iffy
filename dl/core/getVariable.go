package core

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

// GetVariable reads a value of the specified name from the current scope.
// ( ex. loop locals, or -- in a noun scope -- might translate "apple" to "$macintosh" )
type GetVar struct {
	Name rt.TextEval // uses text eval to make template expressions easier
	// if true, and the command is being used to get some text
	// and no variable can be found in the current context of the requested name,
	// see if the requested name is an object instead.
	TryTextAsObject bool `if:"internal"`
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
	if _, p, e := op.getVariableByName(run); e != nil {
		err = cmdError(op, e)
	} else if v, e := p.GetBool(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetNumber(run rt.Runtime) (ret float64, err error) {
	if _, p, e := op.getVariableByName(run); e != nil {
		err = cmdError(op, e)
	} else if v, e := p.GetNumber(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetVar) GetText(run rt.Runtime) (ret string, err error) {
	switch n, p, e := op.getVariableByName(run); e.(type) {
	default:
		err = cmdError(op, e)
	case nil:
		if v, e := p.GetText(run); e != nil {
			err = cmdError(op, e)
		} else {
			ret = v
		}
	case rt.UnknownTarget, rt.UnknownField:
		if !op.TryTextAsObject {
			err = cmdError(op, e)
		} else if id, e := getObjectExactly(run, n); e != nil {
			err = cmdError(op, e)
		} else {
			ret = id
		}
	}
	return
}

// allows us to use GetVar directly in things that take an object.
// in this case, we unbox the variable and assume the text in it is an object id.
func (op *GetVar) GetObjectRef(run rt.Runtime) (retId string, err error) {
	if _, p, e := op.getVariableByName(run); e != nil {
		err = cmdError(op, e)
	} else if str, e := p.GetText(run); e != nil {
		err = cmdError(op, e)
	} else if !strings.HasPrefix(str, "#") {
		e := errutil.New("stored name isnt an object", str)
		err = cmdError(op, e)
	} else {
		retId = str
	}
	return
}

func (op *GetVar) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	// if _, p, e := op.getVariableByName(run); e != nil {
	// 	err = cmdError(op, e)
	// } else if v, e := p.GetNumberStream(run); e != nil {
	// 	err = cmdError(op, e)
	// } else {
	// 	ret = v
	// }
	err = errutil.New("not implemented")
	return
}

func (op *GetVar) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	// if _, p, e := op.getVariableByName(run); e != nil {
	// 	err = cmdError(op, e)
	// } else if v, e := p.GetTextStream(run); e != nil {
	// 	err = cmdError(op, e)
	// } else {
	// 	ret = v
	// }
	err = errutil.New("not implemented")
	return
}

// GetVar asks for a variable using a text eval;
// we first need to determine which actual variable name they mean.
func (op *GetVar) getVariableByName(run rt.Runtime) (retName string, retValue rt.Value, err error) {
	// first resolve the requested variable name into text
	if n, e := rt.GetText(run, op.Name); e != nil {
		err = e
	} else {
		retName = n // then try to get the variable of that name
		retValue, err = run.GetField(object.Variables, n)
	}
	return
}
