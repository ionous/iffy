package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// FromPattern helps runs a pattern
type FromPattern struct {
	Name       string       // i  guess a text eval here would be like a function pointer.
	Parameters []*Parameter // FIX*spec would work if we were an interface
}

type DetermineAct FromPattern
type DetermineNum FromPattern
type DetermineText FromPattern
type DetermineBool FromPattern
type DetermineNumList FromPattern
type DetermineTextList FromPattern

type Parameters []*Parameter

type Parameter struct {
	Name string
	From Assignment
}

func (*Parameter) Compose() composer.Spec {
	return composer.Spec{
		Name:  "parameter",
		Group: "patterns",
	}
}

// Stitch find the pattern, builds the scope, and executes the passed callback to generate a result.
// Its an adapter from the the specific DetermineActivity, DetermineNumber, etc. statements.
func (op *FromPattern) Stitch(run rt.Runtime, fn func(p interface{}) error) (err error) {
	if p, e := run.GetField(op.Name, object.Pattern); e != nil {
		err = e
	} else {
		parms := make(scope.Parameters)
		for _, a := range op.Parameters {
			if e := a.From.Assign(run, func(i interface{}) (err error) {
				return parms.SetVariable(a.Name, i)
			}); e != nil {
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

func (*DetermineAct) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_act",
		Group: "patterns",
	}
}

func (op *DetermineAct) Execute(run rt.Runtime) (err error) {
	err = (*FromPattern)(op).Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.Execute); !ok {
			err = errutil.New("Pattern", op.Name, "not an activity")
		} else {
			err = rt.Run(run, eval)
		}
		return
	})
	return
}

func (*DetermineNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_num",
		Group: "patterns",
	}
}

func (op *DetermineNum) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = (*FromPattern)(op).Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.NumberEval); !ok {
			err = errutil.New("Pattern", op.Name, "not a number")
		} else {
			ret, err = rt.GetNumber(run, eval)
		}
		return
	})
	return
}
func (*DetermineText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_text",
		Group: "patterns",
	}
}
func (op *DetermineText) GetText(run rt.Runtime) (ret string, err error) {
	err = (*FromPattern)(op).Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.TextEval); !ok {
			err = errutil.New("Pattern", op.Name, "not text")
		} else {
			ret, err = rt.GetText(run, eval)
		}
		return
	})
	return
}
func (*DetermineBool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_bool",
		Group: "patterns",
	}
}

func (op *DetermineBool) GetBool(run rt.Runtime) (ret bool, err error) {
	err = (*FromPattern)(op).Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.BoolEval); !ok {
			err = errutil.New("Pattern", op.Name, "not a boolean")
		} else {
			ret, err = rt.GetBool(run, eval)
		}
		return
	})
	return
}

func (*DetermineNumList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_num_list",
		Group: "patterns",
	}
}
func (op *DetermineNumList) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = (*FromPattern)(op).Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.NumListEval); !ok {
			err = errutil.New("Pattern", op.Name, "not a boolean")
		} else {
			ret, err = rt.GetNumberStream(run, eval)
		}
		return
	})
	return
}

func (*DetermineTextList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_text_list",
		Group: "patterns",
	}
}

func (op *DetermineTextList) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = (*FromPattern)(op).Stitch(run, func(p interface{}) (err error) {
		if eval, ok := p.(rt.TextListEval); !ok {
			err = errutil.New("Pattern", op.Name, "not a boolean")
		} else {
			ret, err = rt.GetTextStream(run, eval)
		}
		return
	})
	return
}
