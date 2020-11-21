package core

import (
	"strings"

	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// ObjectName implements ObjectEval, searching for an object named as specified.
// ex. ObjectName{ "target" } looks for the object named "target".
// this is an internal command, used by express.... fix: maybe it should live there.
// and maybe rename to something like "Get/FindObjectId"
type ObjectName struct {
	Name rt.TextEval
}

// tbd: this isnt currently exposed....
// this also could be done by IsKindOf("kind")
type ObjectExists struct {
	Obj rt.ObjectEval
}

// NameOf returns the full name of an object as declared by the author.
// It doesnt change over the course of play. To change the name use the "printed name" property.
type NameOf struct {
	Obj rt.ObjectEval
}

// KindOf returns the class of an object.
type KindOf struct {
	Obj rt.ObjectEval
}

// IsKindOf is less about caring, and more about sharing;
// it returns true when the object is compatible with the named kind.
type IsKindOf struct {
	Obj  rt.ObjectEval
	Kind rt.TextEval
}

// IsExactKindOf returns true when the object is of exactly the named kind.
type IsExactKindOf struct {
	Obj  rt.ObjectEval
	Kind rt.TextEval
}

func (*ObjectName) Compose() composer.Spec {
	return composer.Spec{
		Name:  "object_name",
		Group: "internal",
		Desc:  "Object Name: Returns a noun's object id.",
		Spec:  "object named {name:text_eval}",
	}
}

func (op *ObjectName) GetObject(run rt.Runtime) (ret g.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if v, e := getObjectExactly(run, name.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (*ObjectExists) Compose() composer.Spec {
	return composer.Spec{
		Name:  "object_exists",
		Group: "objects",
		Desc:  "Object Exists: Returns whether there is a noun of the specified name.",
		Spec:  "object named {name:text_eval}",
	}
}

func (op *ObjectExists) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if _, e := safe.GetObject(run, op.Obj); e != nil {
		ret = g.False
	} else {
		ret = g.True
	}
	return

	// fix? b/c of cmdError ( errutil.multierror) we can't test the error type like this.
	// switch _, e := safe.GetObject(run, op.Obj); e.(type) {
	// case nil:
	// 	ret = g.True
	// case g.UnknownObject:
	// 	ret = g.False
	// default:
	// 	err = cmdError(op, e)
	// }
}

// find an object with the passed partial name
func getObjectExactly(run rt.Runtime, name string) (ret g.Value, err error) {
	switch val, e := run.GetField(object.Value, name); e.(type) {
	case g.UnknownField:
		err = g.UnknownObject(name)
	default:
		ret, err = val, e
	}
	return
}

// first look for a variable named "name" in scope, unbox it (if need be) to return the object's id.
func getObjectInexactly(run rt.Runtime, name string) (ret g.Value, err error) {
	switch val, e := safe.Variable(run, name, ""); e.(type) {
	default:
		err = e
	// if there's no such variable, check if there's an object of that name.
	case g.UnknownTarget, g.UnknownField:
		ret, err = getObjectByName(run, name)
	case nil:
		ret = val
	}
	return
}

func (*NameOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "name_of",
		Group: "objects",
		Desc:  "Name Of: Full name of the object.",
		Spec:  "name of {object%obj:object_ref}",
	}
}

func (op *NameOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.Field(run, op.Obj, object.Name, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (*KindOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "kind_of",
		Group: "objects",
		Desc:  "Kind Of: Friendly name of the object's kind.",
		Spec:  "kind of {object%obj:object_ref}",
	}
}

func (op *KindOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.Field(run, op.Obj, object.Kind, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (*IsKindOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "is_kind_of",
		Spec:  "Is {object%obj:object_ref} a kind of {kind:singular_kind}",
		Group: "objects",
		Desc:  "Is Kind Of: True if the object is compatible with the named kind.",
	}
}

func (op *IsKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if tgtKind, e := safe.GetText(run, op.Kind); e != nil {
		err = cmdError(op, e)
	} else if fullPath, e := safe.Field(run, op.Obj, object.Kinds, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		// Contains reports whether second is within first.
		b := strings.Contains(fullPath.String()+",", tgtKind.String()+",")
		ret = g.BoolOf(b)
	}
	return
}

func (*IsExactKindOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "is_exact_class",
		Group: "objects",
		Desc:  "Is Exact Kind: True if the object is exactly the named kind.",
	}
}

func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if tgtKind, e := safe.GetText(run, op.Kind); e != nil {
		err = cmdError(op, e)
	} else if kind, e := safe.Field(run, op.Obj, object.Kind, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		b := tgtKind.String() == kind.String()
		ret = g.BoolOf(b)
	}
	return
}
