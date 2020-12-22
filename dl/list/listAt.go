package list

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type At struct {
	List  string // variable name
	Index rt.NumberEval
}

// future: lists of lists? probably through lists of records containing lists.
func (*At) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_at",
		Group: "list",
		Spec:  "list {list:text} at {index:number_eval}",
		Desc:  "Value of List: Get a value from a list. The first element is is index 1.",
	}
}

func (op *At) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.NumList)
}

func (op *At) GetText(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.TextList)
}

func (op *At) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.getAt(run, affine.RecordList)
}

func (op *At) getAt(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if vs, e := safe.Variable(run, op.List, aff); e != nil {
		err = cmdError(op, e)
	} else if idx, e := safe.GetNumber(run, op.Index); e != nil {
		err = cmdError(op, e)
	} else if i, e := safe.Range(idx.Int()-1, 0, vs.Len()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs.Index(i)
	}
	return
}
