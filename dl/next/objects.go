package next

import (
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

func (op *Exists) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else {
		err = run.GetObject(obj, object.Exists, &ret)
	}
	return
}

func (op *KindOf) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else {
		err = run.GetObject(obj, object.Kind, &ret)
	}
	return
}

func (op *IsKindOf) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if tgtKind, e := rt.GetText(run, op.Kind); e != nil {
		err = e
	} else {
		var objKind string
		if e := run.GetObject(obj, object.Kind, &objKind); e != nil {
			err = e
		} else {
			ret = run.IsCompatible(objKind, tgtKind)
		}
	}
	return
}

func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if tgtKind, e := rt.GetText(run, op.Kind); e != nil {
		err = e
	} else {
		var objKind string
		if e := run.GetObject(obj, object.Kind, &objKind); e != nil {
			err = e
		} else {
			ret = objKind == tgtKind
		}
	}
	return
}
