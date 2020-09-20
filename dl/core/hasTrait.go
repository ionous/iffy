package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

// HasTrait a property value from an object by name.
type HasTrait struct {
	Obj   ObjectRef
	Trait rt.TextEval
}

// should be "When the target is publicly named"
func (*HasTrait) Compose() composer.Spec {
	return composer.Spec{
		Name:  "has_trait",
		Spec:  "{object%obj:object_ref} is {trait:text_eval}",
		Group: "objects",
		Desc:  "Has Trait: Return true if the object current is using the named trait.",
	}
}

func (op *HasTrait) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if trait, e := rt.GetText(run, op.Trait); e != nil {
		err = e
	} else if p, e := run.GetField(obj+"."+trait, object.Aspect); e != nil {
		err = e
	} else if aspect, e := p.GetText(run); e != nil {
		err = e
	} else if currTrait, e := run.GetField(obj, aspect); e != nil {
		err = e
	} else if currTrait, e := currTrait.GetText(run); e != nil {
		ret = trait == currTrait
	}
	return
}
