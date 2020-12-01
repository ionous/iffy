package rel

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

type Relate struct {
	A, B     rt.ObjectEval
	Relation string
}

func (*Relate) Compose() composer.Spec {
	return composer.Spec{
		Name:  "rel_relate",
		Group: "relations",
		Desc:  "Relate: Relate two nouns.",
	}
}

func (op *Relate) Execute(run rt.Runtime) (err error) {
	if e := op.setRelation(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Relate) setRelation(run rt.Runtime) (err error) {
	if a, e := safe.GetObject(run, op.A); e != nil {
		err = e
	} else if b, e := safe.GetObject(run, op.B); e != nil {
		err = e
	} else {
		err = run.RelateTo(a.String(), b.String(), op.Relation)
	}
	return
}
