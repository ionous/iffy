package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

// ChooseAction runs one block of instructions or another based on the results of a conditional evaluation.
type ChooseAction struct {
	If   rt.BoolEval `if:"selector,placeholder=choose"`
	Do   Activity
	Else Brancher `if:"selector,optional"`
}

// Brancher connects else and else-if clauses.
type Brancher interface {
	Branch(rt.Runtime) error
}

// ChooseMore provides an else-if clause.
// Like ChooseAction it chooses a branch based on an if statement.
type ChooseMore ChooseAction

type ChooseNothingElse struct {
	Do Activity `if:"selector"`
}

func (*ChooseAction) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "if", Role: composer.Command},
	}
}

func (*ChooseMore) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "elseIf", Role: composer.Selector},
	}
}

func (*ChooseNothingElse) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "elseDo", Role: composer.Selector},
	}
}

func (op *ChooseAction) Execute(run rt.Runtime) (err error) {
	if e := op.ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseAction) ifDoElse(run rt.Runtime) (err error) {
	if b, e := safe.GetBool(run, op.If); e != nil {
		err = e
	} else if b.Bool() {
		err = op.Do.Execute(run)
	} else if branch := op.Else; branch != nil {
		err = branch.Branch(run)
	}
	return
}

func (op *ChooseMore) Branch(run rt.Runtime) (err error) {
	if e := (*ChooseAction)(op).ifDoElse(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *ChooseNothingElse) Branch(run rt.Runtime) (err error) {
	if e := op.Do.Execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}
