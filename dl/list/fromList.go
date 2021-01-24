package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type FromNumList struct {
	Var core.Variable `if:"selector"`
}
type FromTxtList struct {
	Var core.Variable `if:"selector"`
}
type FromRecList struct {
	Var core.Variable `if:"selector"`
}

type ListSource interface {
	Affinity() affine.Affinity
	GetListSource(run rt.Runtime) (ret g.Value, err error)
}

func GetListSource(run rt.Runtime, src ListSource) (ret g.Value, err error) {
	if src == nil {
		err = errutil.New("missing source list")
	} else {
		ret, err = src.GetListSource(run)
	}
	return
}
func (*FromNumList) Affinity() affine.Affinity { return affine.NumList }

func (*FromNumList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Uses a list of numbers",
	}
}

func (*FromTxtList) Affinity() affine.Affinity { return affine.TextList }

func (*FromTxtList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Uses a list of text",
	}
}

func (*FromRecList) Affinity() affine.Affinity { return affine.RecordList }

func (*FromRecList) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Uses a list of records",
	}
}

func (op *FromNumList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), op.Affinity()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *FromRecList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), op.Affinity()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *FromTxtList) GetListSource(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Var.String(), op.Affinity()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
