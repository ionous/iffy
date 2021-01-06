package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type ChooseAction struct {
	If   rt.BoolEval `if:"selector,placeholder=choose"`
	Do   Activity
	Else Brancher `if:"selector,optional"`
}

type Brancher interface {
	Branch(run rt.Runtime) (err error)
}

type ChooseMore struct {
	If   rt.BoolEval `if:"selector,placeholder=choose"`
	Do   Activity
	Else Brancher `if:"selector,optional"`
}

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
	return
}

func (op *ChooseMore) Branch(run rt.Runtime) (err error) {
	return
}

func (op *ChooseNothingElse) Branch(run rt.Runtime) (err error) {
	return
}
