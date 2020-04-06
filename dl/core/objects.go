package next

import (
	"strings"

	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

// Exists checks for an object by name.
type Exists struct {
	Obj rt.TextEval
}

// KindOf returns the class of an object.
type KindOf struct {
	Obj rt.TextEval
}

// IsKindOf  is less about caring, and more about sharing;
// it returns true when the object is compatible with the named kind.
type IsKindOf struct {
	Obj, Kind rt.TextEval
}

// IsExactKindOf  returns true when the object is of exactly the named kind.
type IsExactKindOf struct {
	Obj, Kind rt.TextEval
}

func (*Exists) Compose() composer.Spec {
	return composer.Spec{
		Name:  "exists",
		Group: "objects",
		Desc:  "Exists: True if the named object exists.",
	}
}

func (op *Exists) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if p, e := run.GetField(obj, object.Exists); e != nil {
		err = e
	} else {
		ret, err = assign.ToBool(p)
	}
	return
}

func (*KindOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "class_name",
		Group: "objects",
		Desc:  "Kind Of: Friendly name of the object's class.",
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
		Name:  "is_class",
		Spec:  "Is $OBJ a kind of $CLASS",
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
