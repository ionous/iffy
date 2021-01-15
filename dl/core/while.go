package core

import (
	"errors"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

type While struct {
	True rt.BoolEval `if:"selector=while"`
	Do   Activity
}

// MaxLoopError provides both an error and a counter
type MaxLoopError int

func (e MaxLoopError) Error() string { return "nearly infinite loop detected" }

var MaxLoopIterations MaxLoopError = 0xbad

func (*While) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "repeating", Role: composer.Command},
		Group:  "flow",
		Desc:   "Repeat a series of statements while a conditional is true.",
	}
}

func (op *While) Execute(run rt.Runtime) (err error) {
	if !op.Do.Empty() {
	LoopBreak:
		for i := 0; ; i++ {
			if i >= int(MaxLoopIterations) {
				err = cmdError(op, MaxLoopIterations)
				break
			} else if keepGoing, e := safe.GetBool(run, op.True); e != nil {
				err = cmdError(op, e)
				break
			} else if !keepGoing.Bool() {
				// all done
				break
			} else {
				// run the loop:
				if e := op.Do.Execute(run); e != nil {
					var i DoInterrupt
					if !errors.As(e, &i) {
						err = cmdError(op, e)
						break LoopBreak
					} else if !i.KeepGoing {
						break LoopBreak
					}
				}
			}
		}
	}
	return
}
