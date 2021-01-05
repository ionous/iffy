package rel

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

type Relate struct {
	Obj, ToObj rt.TextEval
	Via        string // fix: a relation string.
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
	if a, e := safe.ObjectFromText(run, op.Obj); e != nil {
		err = e
	} else if b, e := safe.ObjectFromText(run, op.ToObj); e != nil {
		err = e
	} else {
		err = run.RelateTo(a.String(), b.String(), op.Via)
	}
	return
}
