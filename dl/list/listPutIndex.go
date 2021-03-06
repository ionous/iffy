package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

/**
 * put: eval(num,txt,rec),
 * intoNum/Txt/RecList: varName,
 * atIndex: numEval.
 */
type PutIndex struct {
	From    core.Assignment `if:"selector"`
	Into    ListTarget      `if:"selector"`
	AtIndex rt.NumberEval
}

func (*PutIndex) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Desc:   "Put: replace one value in a list with another",
	}
}

func (op *PutIndex) Execute(run rt.Runtime) (err error) {
	return
}
