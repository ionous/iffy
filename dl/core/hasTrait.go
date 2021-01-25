package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// HasTrait a property value from an object by name.
type HasTrait struct {
	Object rt.TextEval
	Trait  rt.TextEval
}

// should be "When the target is publicly named"
func (*HasTrait) Compose() composer.Spec {
	return composer.Spec{
		Spec:  "{object:text_eval} is {trait:text_eval}",
		Group: "objects",
		Desc:  "Has Trait: Return true if noun is currently in the requested state.",
	}
}

func (op *HasTrait) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj == nil {
		err = cmdError(op, g.NothingObject)
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else if p, e := obj.FieldByName(trait.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = p
	}
	return
}
