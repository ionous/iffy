package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// FromSourceFields: part of GetAtField
type FromSourceFields interface {
	GetSourceFields(run rt.Runtime) (g.Value, error)
}

// Sources an object with a computed name.
type FromObj struct {
	Object rt.TextEval `if:"selector"`
}

// Sources a recorded stored in a record
// RenderRec implements core.FromSourceFields and simply returns the passed record.
// ( This is used in chains of variable names a.b.c.d )
type FromRec struct {
	Rec rt.RecordEval `if:"selector"`
}

// FromVar returns a record or object from a variable.
type FromVar struct {
	Var Variable `if:"selector"`
}

func (*FromObj) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets an object with a computed name.",
	}
}

func (*FromRec) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a record stored in a record.",
	}
}

func (*FromVar) Compose() composer.Spec {
	// FIX? it'd be great to let this "run-in" to the parent
	// "get: <var>" instead of "get var: <var>"
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Targets a record stored in a variable.",
	}
}

// GetSourceFields returns an object supporting field access.
func (op *FromObj) GetSourceFields(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		err = cmdError(op, g.NothingObject)
	} else {
		ret = v
	}
	return
}

// GetSourceFields returns a record supporting field access.
func (op *FromRec) GetSourceFields(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetRecord(run, op.Rec); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

// GetSourceFields returns a record or object supporting field access.
// ( see also IntoVar )
func (op *FromVar) GetSourceFields(run rt.Runtime) (ret g.Value, err error) {
	if v, e := fieldsFromVar(run, op.Var.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func fieldsFromVar(run rt.Runtime, name string) (ret g.Value, err error) {
	if v, e := run.GetField(object.Variables, name); e != nil {
		err = e
	} else {
		switch aff := v.Affinity(); aff {
		case affine.Record:
			ret = v
		case affine.Text:
			ret, err = safe.ObjectFromString(run, v.String())
		default:
			err = errutil.Fmt("unexpected %q for %q", aff, name)
		}
	}
	return
}

func GetSourceFields(run rt.Runtime, fields FromSourceFields) (ret g.Value, err error) {
	if fields == nil {
		err = safe.MissingEval("target fields")
	} else {
		ret, err = fields.GetSourceFields(run)
	}
	return
}
