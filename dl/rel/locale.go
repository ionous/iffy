package rel

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Locale struct {
	Object rt.ObjectEval
}

func (*Locale) Compose() composer.Spec {
	return composer.Spec{
		Name:  "rel_locale",
		Group: "relations",
		Desc:  "Locale: Return the registered parent of an object.",
	}
}

func (op *Locale) GetObject(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getLocale(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *Locale) getLocale(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.GetObject(run, op.Object); e != nil {
		err = e
	} else if v, e := a.FieldByName(object.Locale); e != nil {
		err = e
	} else if va := v.Affinity(); va != affine.Object {
		err = errutil.New("expected object, got %s", va)
	} else {
		ret = v
	}
	return
}
