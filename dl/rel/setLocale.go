package rel

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

type SetLocale struct {
	Object, Parent rt.ObjectEval
}

func (*SetLocale) Compose() composer.Spec {
	return composer.Spec{
		Name:  "rel_set_locale",
		Group: "relations",
		Desc:  "Set Locale: Sets the registered parent of an object.",
	}
}

func (op *SetLocale) Execute(run rt.Runtime) (err error) {
	if e := op.setLocale(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *SetLocale) setLocale(run rt.Runtime) (err error) {
	if child, e := safe.GetObject(run, op.Object); e != nil {
		err = e
	} else if parent, e := safe.GetObject(run, op.Parent); e != nil {
		err = e
	} else {
		err = child.SetFieldByName(object.Locale, parent)
	}
	return
}
