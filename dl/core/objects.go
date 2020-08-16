package core

import (
	"strings"

	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

// ObjectRef checks for an object by name.
// Implementations also generally implement GetText and GetBool
type ObjectRef interface {
	GetObjectRef(run rt.Runtime) (retName string, retExact bool, err error)
	rt.TextEval // see getObjectFullName
	rt.BoolEval // see getObjectExists
}

// ObjectName implements ObjectRef, searching for an object named exactly as specified.
// It matches rt.TextEval, exists for differentiation in the composer.
type ObjectName struct {
	Name    rt.TextEval
	Exactly bool
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

// IsExactKindOf  returns true when the object is of exactly the named kind.
type IsExactKindOf struct {
	Obj  ObjectRef
	Kind rt.TextEval
}

func (*ObjectName) Compose() composer.Spec {
	return composer.Spec{
		Name:  "object_name",
		Group: "objects",
		Desc:  "ObjectName: Returns a noun's full name, can also be used in true/false statements to determine if the named noun exists.",
		Spec:  "the object named {?exactly} {name:text_eval}",
	}
}

func (op *ObjectName) GetText(run rt.Runtime) (ret string, err error) {
	return getObjectFullName(run, op)
}

func (op *ObjectName) GetBool(run rt.Runtime) (ret bool, err error) {
	return getObjectExists(run, op)
}

func (op *ObjectName) GetObjectRef(run rt.Runtime) (ret string, exact bool, err error) {
	if name, e := rt.GetText(run, op.Name); e != nil {
		err = e
	} else {
		ret, exact = name, op.Exactly
	}
	return
}

func getObjectExists(run rt.Runtime, ref ObjectRef) (okay bool, err error) {
	if n, e := getObjectFullName(run, ref); e != nil {
		err = e
	} else {
		okay = len(n) > 0
	}
	return
}

// returns the object's full name
func getObjectFullName(run rt.Runtime, ref ObjectRef) (ret string, err error) {
	if name, exactly, e := ref.GetObjectRef(run); e != nil {
		err = e
	} else if exactly {
		ret, err = getObjectExactly(run, name)
	} else {
		ret, err = getObjectInexactly(run, name)
	}
	return
}

func getObjectExactly(run rt.Runtime, name string) (ret string, err error) {
	if b, e := run.GetField(name, object.Exists); e != nil {
		err = e
	} else if exists, e := assign.ToBool(b); e != nil {
		err = e
	} else if exists {
		ret = name
	}
	return
}

func getObjectInexactly(run rt.Runtime, name string) (ret string, err error) {
	// first look for the name in scope, the top scope (NounScope) will look globally if need be.
	if local, e := run.GetVariable(name); e != nil {
		err = e
	} else if str, e := assign.ToString(local); e != nil {
		err = e
	} else {
		// then, search by inexact name of the unpacked local
		// ( fix? if its from NounScope, then this is a double lookup;
		// maybe NounScope isnt needed now that we have ObjectName )
		if n, e := run.GetField(str, object.Name); e != nil {
			err = e
		} else {
			ret, err = assign.ToString(n)
		}
	}
	return
}

func (*KindOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "kind_of",
		Group: "objects",
		Desc:  "Kind Of: Friendly name of the object's kind.",
		Spec:  "the kind of {object%obj:object_ref}",
	}
}

func (op *KindOf) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if p, e := run.GetField(obj, object.Kind); e != nil {
		err = e
	} else {
		ret, err = assign.ToString(p)
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
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if tgtKind, e := rt.GetText(run, op.Kind); e != nil {
		err = e
	} else if p, e := run.GetField(obj, object.Kinds); e != nil {
		err = e
	} else if fullPath, e := assign.ToString(p); e != nil {
		err = e
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
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if tgtKind, e := rt.GetText(run, op.Kind); e != nil {
		err = e
	} else if p, e := run.GetField(obj, object.Kind); e != nil {
		err = e
	} else if objKind, e := assign.ToString(p); e != nil {
		err = e
	} else {
		ret = objKind == tgtKind
	}
	return
}
