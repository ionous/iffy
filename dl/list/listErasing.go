package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

/**
 * removing: count
 * numFrom/txt/rec:varName,
 * atIndex:
 * as: string, do:{}
 */
type Erasing struct {
	EraseIndex
	As string // fix: variable definition field
	Do core.Activity
}

func (op *Erasing) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erasing", Role: composer.Command},
		Group:  "list",
		Desc: `Erasing from list: Erase elements from the front or back of a list.
Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.`,
	}
}

func (op *Erasing) Execute(run rt.Runtime) (err error) {
	if e := op.popping(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Erasing) popping(run rt.Runtime) (err error) {
	if els, e := op.pop(run); e != nil {
		err = e
	} else {
		run.PushScope(&scope.TargetValue{op.As, els})
		err = op.Do.Execute(run)
		run.PopScope()
	}
	return
}
