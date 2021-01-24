package core

import (
	"strings"

	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type ObjectExists struct {
	Object rt.TextEval `if:"selector=valid,placeholder=object"`
}

// NameOf returns the full name of an object as declared by the author.
// It doesnt change over the course of play. To change the name use the "printed name" property.
type NameOf struct {
	Object rt.TextEval
}

// KindOf returns the class of an object.
type KindOf struct {
	Object rt.TextEval
}

// IsKindOf is less about caring, and more about sharing;
// it returns true when the object is compatible with the named kind.
type IsKindOf struct {
	Object rt.TextEval `if:"selector"`
	Kind   string      `if:"selector=is"`
}

type IsExactKindOf struct {
	Object rt.TextEval `if:"selector"`
	Kind   string      `if:"selector=isExactly"`
}

func (*ObjectExists) Compose() composer.Spec {
	return composer.Spec{
		Group:  "objects",
		Desc:   "Object Exists: Returns whether there is a object of the specified name.",
		Fluent: &composer.Fluid{Name: "is", Role: composer.Function},
	}
}

func (op *ObjectExists) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); obj == nil {
		ret = g.False
	} else if e == nil {
		ret = g.True
	} else if _, isUnknown := e.(g.UnknownObject); isUnknown {
		ret = g.False
	} else {
		err = e
	}
	return
}

func (*NameOf) Compose() composer.Spec {
	return composer.Spec{
		Group: "objects",
		Desc:  "Name Of: Full name of the object.",
		Spec:  "name of {object:text_eval}",
	}
}

func (op *NameOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj == nil {
		ret = g.StringOf("") // fix: or, should it be "nothing"
	} else if v, e := safe.Unpack(obj, object.Name, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (*KindOf) Compose() composer.Spec {
	return composer.Spec{
		Group: "objects",
		Desc:  "Kind Of: Friendly name of the object's kind.",
		Spec:  "kind of {object:text_eval}",
	}
}

func (op *KindOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj == nil {
		ret = g.StringOf("")
	} else if v, e := safe.Unpack(obj, object.Kind, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (*IsKindOf) Compose() composer.Spec {
	return composer.Spec{
		Group:  "objects",
		Fluent: &composer.Fluid{Name: "kindOf", Role: composer.Function},
		Desc:   "Is Kind Of: True if the object is compatible with the named kind.",
	}
}

func (op *IsKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj == nil {
		ret = g.BoolOf(false)
	} else if fullPath, e := safe.Unpack(obj, object.Kinds, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		// Contains reports whether second is within first.
		kind := lang.Breakcase(op.Kind)
		b := strings.Contains(fullPath.String()+",", kind+",")
		ret = g.BoolOf(b)
	}
	return
}

func (*IsExactKindOf) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "kindOf", Role: composer.Function},
		Group:  "objects",
		Desc:   "Is Kind Of: True if the object is compatible with the named kind.",
	}
}

func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj == nil {
		ret = g.BoolOf(false)
	} else if fullPath, e := safe.Unpack(obj, object.Kinds, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		// Contains reports whether second is within first.
		kind := lang.Breakcase(op.Kind)
		b := strings.Contains(fullPath.String()+",", kind+",")
		ret = g.BoolOf(b)
	}
	return
}
