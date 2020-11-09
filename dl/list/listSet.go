package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type Set struct {
	List    string // variable name
	Index   rt.NumberEval
	Element core.Assignment
}

func (*Set) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_set",
		Group: "list",
		// Spec:  "{list:text} entry {index:number}",
		Desc: "Set Value of List: Overwrite an existing value in a list.",
	}
}

func (op *Set) Execute(run rt.Runtime) (err error) {
	if e := op.setAt(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Set) setAt(run rt.Runtime) (err error) {
	if els, e := run.GetField(object.Variables, op.List); e != nil {
		err = e
	} else if i, e := rt.GetNumber(run, op.Index); e != nil {
		err = e
	} else if cnt, e := els.GetLen(); e != nil {
		err = e
	} else if i := int(i - 1); i < 0 || i >= cnt {
		err = rt.OutOfRange{i, cnt}
	} else if els, e := op.replace(run, els, i); e != nil {
		err = e
	} else {
		err = run.SetField(object.Variables, op.List, els)
	}
	return
}

func (op *Set) replace(run rt.Runtime, els g.Value, i int) (ret g.Value, err error) {
	switch a := els.Affinity(); a {
	case affine.NumList:
		if els, e := els.GetNumList(); e != nil {
			err = e
		} else if el, e := core.GetAssignedValue(run, op.Element); e != nil {
			err = e
		} else if el, e := el.GetNumber(); e != nil {
			err = e
		} else {
			els[i] = el
			ret = g.FloatsOf(els)
		}
	case affine.TextList:
		if els, e := els.GetTextList(); e != nil {
			err = e
		} else if el, e := core.GetAssignedValue(run, op.Element); e != nil {
			err = e
		} else if el, e := el.GetText(); e != nil {
			err = e
		} else {
			els[i] = el
			ret = g.StringsOf(els)
		}
	default:
		err = errutil.Fmt("variable '%s(%s)' isn't a list", op.List, a)
	}
	return
}
