package core

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Fields interface {
	GetFields(run rt.Runtime) (g.Value, error)
}

// Targets a recorded stored in a variable.
type IntoRec struct {
	Var Variable `if:"selector"`
}

// Targets an object stored in a variable.
type IntoObj struct {
	Var Variable `if:"selector"`
}

// Targets an object with a computed name.
type IntoObjNamed struct {
	ObjName rt.TextEval `if:"selector"`
}

func (*IntoRec) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a record stored in a variable",
	}
}

func (*IntoObj) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets an object with a predetermined name",
	}
}

func (*IntoObjNamed) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets an object with a computed name",
	}
}

func (op *IntoRec) GetFields(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.Variable(run, op.Var.String(), affine.Record); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return

}
func (op *IntoObj) GetFields(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.Variable(run, op.Var.String(), affine.Object); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IntoObjNamed) GetFields(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.ObjName); e != nil {
		err = cmdError(op, e)
	} else {
		ret = obj
	}
	return
}

func GetFields(run rt.Runtime, fields Fields) (ret g.Value, err error) {
	if fields == nil {
		err = safe.MissingEval("empty fields")
	} else {
		ret, err = fields.GetFields(run)
	}
	return
}
