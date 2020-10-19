package core

import (
	"strings"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

// ObjectRef finds an object by name, it returns an id.
type ObjectRef interface {
	// returns UnknownObject when successfully determining there is no such object
	GetObjectRef(run rt.Runtime) (retId string, err error)
}

// ObjectName implements ObjectRef, searching for an object named as specified.
// ex. ObjectName{ "target" }  looks for the object named "target".
// this is an internal command, used by express.... fix: maybe it should live there.
// and maybe rename to something like "Get/FindObjectId"
type ObjectName struct {
	Name rt.TextEval
}

// tbd: this isnt currently exposed....
// this also could be done by IsKindOf("kind")
type ObjectExists struct {
	Obj ObjectRef
}

// NameOf returns the full name of an object as written by the author when declared.
// The name cannot be changed at runtime, instead use the "printed name" property.
type NameOf struct {
	Obj ObjectRef
}

// KindOf returns the class of an object.
type KindOf struct {
	Obj ObjectRef
}

// IsKindOf is less about caring, and more about sharing;
// it returns true when the object is compatible with the named kind.
type IsKindOf struct {
	Obj  ObjectRef
	Kind rt.TextEval
}

// IsExactKindOf returns true when the object is of exactly the named kind.
type IsExactKindOf struct {
	Obj  ObjectRef
	Kind rt.TextEval
}

func GetObjectRef(run rt.Runtime, ref ObjectRef) (retId string, err error) {
	if ref == nil {
		err = rt.MissingEval("empty object ref")
	} else if id, e := ref.GetObjectRef(run); e != nil {
		err = e
	} else {
		retId = id
	}
	return
}

func (*ObjectName) Compose() composer.Spec {
	return composer.Spec{
		Name:  "object_name",
		Group: "internal",
		Desc:  "Object Name: Returns a noun's object id.",
		Spec:  "object named {name:text_eval}",
	}
}

// can be used as text, returns the object id.
func (op *ObjectName) GetText(run rt.Runtime) (ret string, err error) {
	return op.GetObjectRef(run)
}

func (op *ObjectName) GetObjectRef(run rt.Runtime) (retId string, err error) {
	if name, e := rt.GetText(run, op.Name); e != nil {
		err = e
	} else {
		retId, err = getObjectExactly(run, name)
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

func (op *ObjectExists) GetBool(run rt.Runtime) (okay bool, err error) {
	// checking for object.Exists only searches by object id
	// we want to check for the object by friendly name, and possibly by looking in scope
	switch _, e := GetObjectRef(run, op.Obj); e.(type) {
	case nil:
		okay = true
	case rt.UnknownObject:
		okay = false
	default:
		err = cmdError(op, e)
	}
	return
}

// find an object with the passed partial name; return its id
func getObjectExactly(run rt.Runtime, name string) (retId string, err error) {
	if strings.HasPrefix(name, "#") {
		retId = name // ids start with prefix #
	} else {
		switch id, e := run.GetField(object.Id, name); e.(type) {
		case rt.UnknownField:
			err = rt.UnknownObject(name)
		default:
			err = e
		case nil:
			retId, err = id.GetText()
		}
	}
	return
}

// first look for a variable named "name" in scope, unbox it (if need be) to return the object's id.
// this differs from GetVariableOrObject.
// that function tries to find the value of a variable or object id named "something"
// this tries to resolve the name "something" into an object.
func getObjectInexactly(run rt.Runtime, name string) (retId string, err error) {
	switch p, e := run.GetField(object.Variables, name); e.(type) {
	// if there's no such variable, the inexact search then checks if there's an object of that name.
	case rt.UnknownTarget, rt.UnknownField:
		retId, err = getObjectExactly(run, name)
	case nil:
		// if we found such a variable, get its contents and look up the referenced object.
		if unboxedName, e := p.GetText(); e != nil {
			err = e
		} else {
			retId, err = getObjectExactly(run, unboxedName)
		}
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

func (op *NameOf) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := GetObjectRef(run, op.Obj); e != nil {
		err = cmdError(op, e)
	} else if p, e := run.GetField(object.Name, obj); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = p.GetText()
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

func (op *KindOf) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := GetObjectRef(run, op.Obj); e != nil {
		err = cmdError(op, e)
	} else if p, e := run.GetField(object.Kind, obj); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = p.GetText()
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

func (op *IsKindOf) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := GetObjectRef(run, op.Obj); e != nil {
		err = cmdError(op, e)
	} else if tgtKind, e := rt.GetText(run, op.Kind); e != nil {
		err = cmdError(op, e)
	} else if p, e := run.GetField(object.Kinds, obj); e != nil {
		err = cmdError(op, e)
	} else if fullPath, e := p.GetText(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = strings.Contains(fullPath+",", tgtKind+",")
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

func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := GetObjectRef(run, op.Obj); e != nil {
		err = cmdError(op, e)
	} else if tgtKind, e := rt.GetText(run, op.Kind); e != nil {
		err = cmdError(op, e)
	} else if p, e := run.GetField(object.Kind, obj); e != nil {
		err = cmdError(op, e)
	} else if objKind, e := p.GetText(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = objKind == tgtKind
	}
	return
}
