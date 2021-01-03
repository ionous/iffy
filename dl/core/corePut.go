package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

/**
 * put: eval,
 * intoObj/intoObjNamed/intoRec: varName,objName,textEval,
 * atField: string.
 */
type PutAtField struct {
	From    Assignment `if:"unlabeled"`
	Into    Fields     `if:"unlabeled"`
	AtField string
}

func (op *PutAtField) Execute(run rt.Runtime) (err error) {
	return
}

type Fields interface {
	GetFields(run rt.Runtime) (g.Value, error)
}

type IntoObj struct {
	ObjName string `if:"unlabeled"`
}
type IntoObjNamed struct {
	ObjName rt.TextEval `if:"unlabeled"`
}
type IntoRec struct {
	Var Variable `if:"unlabeled"`
}

func (op *IntoObj) GetFields(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *IntoRec) GetFields(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *IntoObjNamed) GetFields(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (*PutAtField) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Desc:   "Put: put a value into the field of an record or object",
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

func (*IntoRec) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a record stored in a variable",
	}
}
