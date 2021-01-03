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
type Erase struct {
	Count   rt.NumberEval  `if:"unlabeled"`
	From    FromListSource `if:"unlabeled"`
	AtIndex rt.NumberEval
}

func (*Erase) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erase", Role: composer.Command},
		Desc:   "Erase: remove one or more values from a list",
	}
}

func (op *Erase) Execute(run rt.Runtime) (err error) {
	return
}
