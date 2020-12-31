package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

/**
 * put: eval(num,txt,rec),
 * intoNum/Txt/RecList: varName,
 * atBack|atFront.
 */
type PutAtEdge struct {
	From   core.Assignment `if:"unlabeled"`
	Into   ListVar         `if:"unlabeled"`
	AtEdge Edge            `if:"unlabeled"`
}

/**
 * put: eval(num,txt,rec),
 * intoNum/Txt/RecList: varName,
 * atIndex: numEval.
 */
type PutAtIndex struct {
	From    core.Assignment `if:"unlabeled"`
	Into    ListVar         `if:"unlabeled"`
	AtIndex rt.NumberEval
}

type IntoNumList struct {
	VarName string `if:"unlabeled"`
}
type IntoTxtList struct {
	VarName string `if:"unlabeled"`
}
type IntoRecList struct {
	VarName string `if:"unlabeled"`
}

func (op *PutAtEdge) Execute(run rt.Runtime) (err error) {
	return
}

func (op *PutAtIndex) Execute(run rt.Runtime) (err error) {
	return
}

type ListVar interface {
	GetListVar(run rt.Runtime) (g.Value, error)
}

func (op *IntoNumList) GetListVar(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *IntoRecList) GetListVar(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *IntoTxtList) GetListVar(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (*PutAtEdge) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Desc:   "Put: add a value to a list",
	}
}

func (*PutAtIndex) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Desc:   "Put: replace one value in a list with another",
	}
}

func (*IntoNumList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of numbers",
	}
}

func (*IntoTxtList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of text",
	}
}

func (*IntoRecList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of records",
	}
}
