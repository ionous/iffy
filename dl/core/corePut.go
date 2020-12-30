package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type PutAtField struct {
	From    Assignment
	Into    Fields
	AtField string
}

func (op *PutAtField) Execute(run rt.Runtime) (err error) {
	return
}

type Fields interface {
	GetFields(run rt.Runtime) (g.Value, error)
}

type IntoObj struct{ ObjName string }

func (op *IntoObj) GetFields(run rt.Runtime) (ret g.Value, err error) {
	return
}

type IntoObjNamed struct{ ObjName rt.TextEval }

type IntoRec struct{ VarName string }

func (op *IntoRec) GetFields(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (op *IntoObjNamed) GetFields(run rt.Runtime) (ret g.Value, err error) {
	return
}

func (*PutAtField) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluency{Name: "put", RunIn: true},
		// Spec:   "put: {%from:assignment} {%into:fields} atField: {%atField:text}",
		Desc: "Put: put a value into the field of an record or object",
	}
}

func (*IntoObj) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluency{RunIn: true},
		Desc:   "Targets an object with a predetermined name",
	}
}

func (*IntoObjNamed) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluency{RunIn: true},
		Desc:   "Targets an object with a computed name",
	}
}

func (*IntoRec) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluency{RunIn: true},
		Desc:   "Targets a record stored in a variable",
	}
}
