package rel

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Relatives struct {
	Object   rt.ObjectEval
	Relation string
}

func (*Relatives) Compose() composer.Spec {
	return composer.Spec{
		Name:  "rel_relatives",
		Group: "relations",
		Desc:  "Relatives: Returns the relatives of a noun as a list of names.",
	}
}

func (op *Relatives) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if vs, e := op.relatives(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.StringsOf(vs)
	}
	return
}

func (op *Relatives) relatives(run rt.Runtime) (ret []string, err error) {
	if a, e := safe.GetObject(run, op.Object); e != nil {
		err = e
	} else {
		ret, err = run.RelativesOf(a.String(), op.Relation)
	}
	return
}
