package render

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// RenderRef returns the value of a variable or the id of an object.
type RenderRef struct {
	core.Var
	Flags TryAsNoun
}

// Compose implements composer.Composer
func (*RenderRef) Compose() composer.Spec {
	return composer.Spec{
		Group: "internal",
	}
}

func (op *RenderRef) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getAssignedValue(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

// GetText handles unpacking a text variable,
// turning an object variable into an id, or
// looking for an object of the passed name ( if no variable of the name exists. )
func (op *RenderRef) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getText(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) getText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getAssignedValue(run); e != nil {
		err = e
	} else if aff := v.Affinity(); aff == affine.Text {
		ret = v
	} else if aff == affine.Object {
		ret = g.ObjectAsText(v)
	} else {
		err = errutil.Fmt("unexpected %q", aff)
	}
	return
}

func (op *RenderRef) getAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := getVariable(run, op.Name, op.Flags); e != nil {
		err = e
	} else if val != nil {
		ret = val
	} else if !op.Flags.tryObject() {
		err = g.UnknownVariable(op.Name)
	} else if obj, e := safe.ObjectFromString(run, op.Name); e != nil {
		err = e
	} else {
		ret = obj
	}
	return
}

// returns nil if the named variable doesnt exist; errors only on critical errors.
func getVariable(run rt.Runtime, n string, flags TryAsNoun) (ret g.Value, err error) {
	if flags.tryVariable() {
		switch v, e := safe.CheckVariable(run, n, ""); e.(type) {
		default:
			err = e
		case g.UnknownTarget, g.UnknownField:
			ret = nil // not a variable
		case nil:
			ret = v
		}
	}
	return
}
