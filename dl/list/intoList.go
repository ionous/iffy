package list

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type ListTarget interface {
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
	if v, e := safe.Variable(run, op.Var.String(), affine.NumList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IntoRecList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.Variable(run, op.Var.String(), affine.RecordList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IntoTxtList) GetListTarget(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.Variable(run, op.Var.String(), affine.TextList); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
