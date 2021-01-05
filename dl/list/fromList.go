package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type ListSource interface {
	GetListSource(run rt.Runtime) (g.Value, error)
}

type FromNumList struct {
	Var core.Variable `if:"unlabeled"`
}
type FromTxtList struct {
	Var core.Variable `if:"unlabeled"`
}
type FromRecList struct {
	Var core.Variable `if:"unlabeled"`
}

func (op *FromNumList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *FromRecList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *FromTxtList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (*FromNumList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of numbers",
	}
}

func (*FromTxtList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of text",
	}
}

func (*FromRecList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a list of records",
	}
}
