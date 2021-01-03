package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

/**
 * put: eval(num,txt,rec),
 * intoNum/Txt/RecList: varName,
 * atBack|atFront.
 */
type PutAtEdge struct {
	From   core.Assignment `if:"unlabeled"`
	Into   IntoListTarget  `if:"unlabeled"`
	AtEdge Edge            `if:"unlabeled"`
}

func (*PutAtEdge) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Desc:   "Put: add a value to a list",
	}
}

func (op *PutAtEdge) Execute(run rt.Runtime) (err error) {
	return
}

/**
 * put: eval(num,txt,rec),
 * intoNum/Txt/RecList: varName,
 * atIndex: numEval.
 */
type PutAtIndex struct {
	From    core.Assignment `if:"unlabeled"`
	Into    IntoListTarget  `if:"unlabeled"`
	AtIndex rt.NumberEval
}

func (*PutAtIndex) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Desc:   "Put: replace one value in a list with another",
	}
}

func (op *PutAtIndex) Execute(run rt.Runtime) (err error) {
	return
}
