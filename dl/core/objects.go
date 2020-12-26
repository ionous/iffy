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
// and maybe rename to something like "Get/FindObjectId"
type ObjectName struct {
	Name rt.TextEval
}

// tbd: this isnt currently exposed....
// this also could be done by IsKindOf("kind")
type ObjectExists struct {
	Object rt.ObjectEval
}

// NameOf returns the full name of an object as declared by the author.
// It doesnt change over the course of play. To change the name use the "printed name" property.
type NameOf struct {
	Object rt.ObjectEval
}

// KindOf returns the class of an object.
type KindOf struct {
	Object rt.ObjectEval
}

// IsKindOf is less about caring, and more about sharing;
// it returns true when the object is compatible with the named kind.
type IsKindOf struct {
	Object rt.ObjectEval
	Kind   string
}

func (*ObjectName) Compose() composer.Spec {
	return composer.Spec{
		Name:  "object_name",
		Group: "objects",
		Desc:  "Object Name: Returns a noun's object id.",
		Spec:  "object named {name:text_eval}",
	}
}

func (op *ObjectName) GetObject(run rt.Runtime) (ret g.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if v, e := getObjectNamed(run, name.String()); e != nil {
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
	if _, e := safe.GetObject(run, op.Object); e != nil {
		ret = g.False
	} else {
		ret = g.True
	}
	return

	// fix? b/c of cmdError ( errutil.multierror) we can't test the error type like this.
	// switch _, e := safe.GetObject(run, op.Object); e.(type) {
	// case nil:
	// 	ret = g.True
	// case g.UnknownObject:
	// 	ret = g.False
	// default:
	// 	err = cmdError(op, e)
	// }
}

func (*NameOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "name_of",
		Group: "objects",
		Desc:  "Name Of: Full name of the object.",
		Spec:  "name of {object:object_eval}",
	}
}

func (op *NameOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.Field(run, op.Object, object.Name, affine.Text); e != nil {
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
		Spec:  "kind of {object:object_eval}",
	}
}

func (op *KindOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.Field(run, op.Object, object.Kind, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (*IsKindOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "is_kind_of",
		Spec:  "is {object:object_eval} a kind of {kind:singular_kind}",
		Group: "objects",
		Desc:  "Is Kind Of: True if the object is compatible with the named kind.",
	}
}

func (op *IsKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if fullPath, e := safe.Field(run, op.Object, object.Kinds, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		// Contains reports whether second is within first.
		b := strings.Contains(fullPath.String()+",", op.Kind+",")
		ret = g.BoolOf(b)
	}
	return
}
