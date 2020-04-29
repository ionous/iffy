package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// FromPattern runs a pattern
type FromPattern struct {
	Name       string // i  guess a text eval here would be like a function pointer.
	Parameters []*Assign
}

type DetermineAct struct {
	*FromPattern
}

type DetermineNum struct {
	*FromPattern
}

type DetermineText struct {
	*FromPattern
}

type DetermineBool struct {
	*FromPattern
}

type DetermineNumList struct {
	*FromPattern
}

type DetermineTextList struct {
	*FromPattern
}

// func (*FromPattern) Compose() composer.Spec {
// 	return composer.Spec{
// 		Name:  "set_num_list_var",
// 		Group: "internal",
// 		Desc:  "Set Number List Variable: Sets the named variable to the passed number list.",
// 	}
// }

// Stitch find the pattern, builds the scope, and executes the passed callback to generate a result.
// Its an adapter from the the specific DetermineActivity, DetermineNumber, etc. statements.
func (op *FromPattern) Stitch(run rt.Runtime, fn func(p interface{}) error) (err error) {
	if p, e := run.GetField(op.Name, object.Pattern); e != nil {
		err = e
	} else {
		parms := make(scope.Parameters)
		for _, a := range op.Parameters {
			if e := a.Assignment.Assign(run, a.Name, &parms); e != nil {
				err = errutil.Append(err, e)
			}
		}
		if err == nil {
			run.PushScope(&parms)
			err = fn(p)
			run.PopScope()
		}
	}
	return
}

func (op *DetermineAct) Execute(run rt.Runtime) (err error) {
	err = op.FromPattern.Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.Execute); !ok {
			err = errutil.New("Pattern", op.Name, "not an activity")
		} else {
			err = rt.Run(run, eval)
		}
		return
	})
	return
}

func (op *DetermineNum) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = op.FromPattern.Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.NumberEval); !ok {
			err = errutil.New("Pattern", op.Name, "not a number")
		} else {
			ret, err = rt.GetNumber(run, eval)
		}
		return
	})
	return
}

func (op *DetermineText) GetText(run rt.Runtime) (ret string, err error) {
	err = op.FromPattern.Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.TextEval); !ok {
			err = errutil.New("Pattern", op.Name, "not text")
		} else {
			ret, err = rt.GetText(run, eval)
		}
		return
	})
	return
}

func (op *DetermineBool) GetBool(run rt.Runtime) (ret bool, err error) {
	err = op.FromPattern.Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.BoolEval); !ok {
			err = errutil.New("Pattern", op.Name, "not a boolean")
		} else {
			ret, err = rt.GetBool(run, eval)
		}
		return
	})
	return
}

func (op *DetermineNumList) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = op.FromPattern.Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.NumListEval); !ok {
			err = errutil.New("Pattern", op.Name, "not a boolean")
		} else {
			ret, err = rt.GetNumberStream(run, eval)
		}
		return
	})
	return
}

func (op *DetermineTextList) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = op.FromPattern.Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.TextListEval); !ok {
			err = errutil.New("Pattern", op.Name, "not a boolean")
		} else {
			ret, err = rt.GetTextStream(run, eval)
		}
		return
	})
	return
}
