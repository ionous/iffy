package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type IntoListTarget interface {
	GetListTarget(run rt.Runtime) (g.Value, error)
}

type IntoNumList struct {
	Var core.Variable `if:"unlabeled"`
}
type IntoTxtList struct {
	Var core.Variable `if:"unlabeled"`
}
type IntoRecList struct {
	Var core.Variable `if:"unlabeled"`
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

func (op *IntoNumList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *IntoRecList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *IntoTxtList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	return
}
