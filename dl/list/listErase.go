package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

/**
 * erase: numEval
 * from: varName,
 * atIndex: num,
 */
type EraseAtIndex struct {
	Count   rt.NumberEval `if:"unlabeled"`
	From    ListSource    `if:"unlabeled"`
	AtIndex rt.NumberEval
}

func (*EraseAtIndex) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erase", Role: composer.Command},
		Desc:   "Erase: remove one or more values from a list",
	}
}

func (op *EraseAtIndex) Execute(run rt.Runtime) (err error) {
	return
}

/**
 * erase: numEval
 * from: varName,
 * atIndex: num,
 */
type EraseAtEdge struct {
	From   ListSource `if:"unlabeled"`
	AtEdge Edge       `if:"unlabeled"`
}

func (*EraseAtEdge) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erase", Role: composer.Command},
		Desc:   "Erase: remove one or more values from a list",
	}
}

func (op *EraseAtEdge) Execute(run rt.Runtime) (err error) {
	return
}
