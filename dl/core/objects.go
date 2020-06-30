package core

import (
	"strings"

	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

// ObjectName checks for an object by name.
type ObjectName struct {
	Name    rt.TextEval
	Exactly bool
}

// KindOf returns the class of an object.
type KindOf struct {
	Obj *ObjectName
}

// IsKindOf  is less about caring, and more about sharing;
// it returns true when the object is compatible with the named kind.
type IsKindOf struct {
	Obj  *ObjectName
	Kind rt.TextEval
}

// IsExactKindOf  returns true when the object is of exactly the named kind.
type IsExactKindOf struct {
	Obj  *ObjectName
	Kind rt.TextEval
}

func (*ObjectName) Compose() composer.Spec {
	return composer.Spec{
		Name:  "object_name",
		Group: "objects",
		Desc:  "ObjectName: Returns a noun's full name, can also be used in true/false statements to determine if the named noun exists.",
		Spec:  "the {global?global_only} object {named:text}",
	}
}

func (op *ObjectName) GetText(run rt.Runtime) (ret string, err error) {
	if n, e := op.getFullName(run); e != nil {
		err = e
	} else {
		ret = n
	}
	return
}

func (op *ObjectName) GetBool(run rt.Runtime) (ret bool, err error) {
	if n, e := op.getFullName(run); e != nil {
		err = e
	} else {
		ret = len(n) > 0
	}
	return
}

// returns the object's full name
func (op *ObjectName) getFullName(run rt.Runtime) (ret string, err error) {
	if name, e := rt.GetText(run, op.Name); e != nil {
		err = e
	} else if op.Exactly {
		ret, err = op.getExactly(run, name)
	} else {
		ret, err = op.getInexactly(run, name)
	}
	return
}

func (op *ObjectName) getExactly(run rt.Runtime, name string) (ret string, err error) {
	if b, e := run.GetField(name, object.Exists); e != nil {
		err = e
	} else if exists, e := assign.ToBool(b); e != nil {
		err = e
	} else if exists {
		ret = name
	}
	return
}

func (op *ObjectName) getInexactly(run rt.Runtime, name string) (ret string, err error) {
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
		Spec:  "the kind of {obj:object_name}",
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
		Spec:  "Is {noun%obj} a kind of {kind:singular_kind}",
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
