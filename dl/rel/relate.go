package rel

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Relate struct {
	Object   rt.TextEval `if:"selector"`
	ToObject rt.TextEval `if:"selector=to"`
	Via      Relation
}

func (*Relate) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "relate", Role: composer.Command},
		Group:  "relations",
		Desc:   "Relate: Relate two nouns.",
	}
}

func (op *Relate) Execute(run rt.Runtime) (err error) {
	if e := op.setRelation(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Relate) setRelation(run rt.Runtime) (err error) {
	if a, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = e
	} else if b, e := safe.ObjectFromText(run, op.ToObject); e != nil {
		err = e
	} else {
		err = run.RelateTo(objectString(a), objectString(b), op.Via.String())
	}
	return
}

func objectString(obj g.Value) (ret string) {
	if obj != nil {
		ret = obj.String()
	}
	return
}
