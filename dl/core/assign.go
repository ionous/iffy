package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// Assignment helps limit variable and parameter assignment to particular contexts
type Assignment interface {
	Assign(rt.Runtime, func(interface{}) error) error
}

// Assign turns an Assignment a normal statement.
type Assign struct {
	Name string // text eval, or string?
	From Assignment
}

type FromBool struct {
	Val rt.BoolEval
}

type FromNum struct {
	Val rt.NumberEval
}

type FromText struct {
	Val rt.TextEval
}

type FromNumList struct {
	Vals rt.NumListEval
}

type FromTextList struct {
	Vals rt.TextListEval
}

func (op *Assign) Execute(run rt.Runtime) (err error) {
	if assign := op.From; assign == nil {
		err = rt.MissingEval("empty assignment")
	} else {
		err = assign.Assign(run, func(i interface{}) error {
			return run.SetVariable(op.Name, i)
		})
	}
	return
}

func (*Assign) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign",
		Spec:  "let {name} be {assignment}",
		Group: "variables",
		Desc:  "Assignment: Sets a variable to a value.",
	}
}

func (*FromBool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_bool",
		Group: "variables",
		Desc:  "Assign Boolean: Assigns the passed boolean value.",
	}
}

func (op *FromBool) Assign(run rt.Runtime, fn func(interface{}) error) (err error) {
	if val, e := rt.GetBool(run, op.Val); e != nil {
		err = e
	} else {
		err = fn(val)
	}
	return
}

func (*FromNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_num",
		Group: "variables",
		Desc:  "Assign Number: Assigns the passed number.",
	}
}

func (op *FromNum) Assign(run rt.Runtime, fn func(interface{}) error) (err error) {
	if val, e := rt.GetNumber(run, op.Val); e != nil {
		err = e
	} else {
		err = fn(val)
	}
	return
}

func (*FromText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_text",
		Group: "variables",
		Desc:  "Assign Text: Assigns the passed piece of text.",
	}
}

func (op *FromText) Assign(run rt.Runtime, fn func(interface{}) error) (err error) {
	if val, e := rt.GetText(run, op.Val); e != nil {
		err = e
	} else {
		err = fn(val)
	}
	return
}

func (*FromNumList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_num_list",
		Group: "variables",
		Desc:  "Assign Number List: Assigns the passed number list.",
	}
}

func (op *FromNumList) Assign(run rt.Runtime, fn func(interface{}) error) (err error) {
	if vals, e := rt.GetNumList(run, op.Vals); e != nil {
		err = e
	} else {
		err = fn(vals)
	}
	return
}

func (*FromTextList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_text_list",
		Group: "variables",
		Desc:  "Assign Text List: Assigns the passed text list.",
	}
}

func (op *FromTextList) Assign(run rt.Runtime, fn func(interface{}) error) (err error) {
	if vals, e := rt.GetTextList(run, op.Vals); e != nil {
		err = e
	} else {
		err = fn(vals)
	}
	return
}
