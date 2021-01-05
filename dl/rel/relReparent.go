package rel

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

type Reparent struct {
	ChildObj, ToParentObj rt.TextEval
}

func (*Reparent) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "reparent", Role: composer.Command},
		Group:  "relations",
		Desc:   "Set Locale: Sets the registered parent of an object.",
	}
}

func (op *Reparent) Execute(run rt.Runtime) (err error) {
	if e := op.setLocale(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Reparent) setLocale(run rt.Runtime) (err error) {
	if child, e := safe.ObjectFromText(run, op.ChildObj); e != nil {
		err = e
	} else if parent, e := safe.ObjectFromText(run, op.ToParentObj); e != nil {
		err = e
	} else {
		err = child.SetFieldByName(object.Locale, parent)
	}
	return
}
