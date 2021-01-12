package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

type While struct {
	True rt.BoolEval `if:"selector=while"`
	Do   Activity
}

func (*While) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "repeating", Role: composer.Command},
		Group:  "flow",
		Desc:   "Repeat a series of statements while a conditional is true.",
	}
}

func (op *While) Execute(run rt.Runtime) (err error) {
	if !op.Do.Empty() {
		for i := 0; i < 10000; i++ {
			if ok, e := safe.GetBool(run, op.True); e != nil || !ok.Bool() {
				err = cmdError(op, e)
				break
			} else if e := op.Do.Execute(run); e != nil {
				err = cmdError(op, e)
				break
			}
		}
	}
	return
}
