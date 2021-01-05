package core

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// Var reads the value of the specified name from the current scope.
type Var struct {
	Name  string
	Flags TryAsNoun `if:"internal"`
}

// Compose implements composer.Composer
func (*Var) Compose() composer.Spec {
	return composer.Spec{
		Name:  "get_var",
		Spec:  "var: {name:text}",
		Group: "variables",
		Desc:  "Get Variable: Return the value of the named variable.",
	}
}

func (op *Var) Affinity() affine.Affinity { return "" }

// GetAssignedValue implements Assignment so we can SetXXX values from variables without a FromXXX statement in between.
func (op *Var) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, "")
}

func (op *Var) GetBool(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Bool)
}

func (op *Var) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Number)
}

func (op *Var) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Record)
}

func (op *Var) GetText(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Text)
}

// allows us to use Var directly in commands which take a named object.
func (op *Var) GetObject(run rt.Runtime) (ret g.Value, err error) {
	if val, e := op.getObject(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (op *Var) getObject(run rt.Runtime) (ret g.Value, err error) {
	// try to get the variable named "name"
	if val, e := getVariableValue(run, op.Name, affine.Object, op.Flags); e != nil {
		err = e
	} else if val != nil {
		ret = val
	} else if !op.Flags.tryObject() {
		err = g.UnknownObject(op.Name)
	} else {
		// try to get the object named "name"
		// this is for internal use by template rendering;
		// fix? maybe that should get its own command
		ret, err = safe.ObjectFromString(run, op.Name)
	}
	return
}

func (op *Var) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.NumList)
}

func (op *Var) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.TextList)
}

// fix: should we bother to try to confirm that it's a RecordList or let the caller figure it out?
// see also: GetObject
func (op *Var) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.RecordList)
}

func (op *Var) getVar(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := safe.Variable(run, op.Name, aff); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

// returns a nil value if the variable couldnt be found
// returns error only critical errors
func getVariableValue(run rt.Runtime, text string, aff affine.Affinity, flags TryAsNoun) (ret g.Value, err error) {
	// first resolve the requested variable name into text
	if flags.tryVariable() {
		switch val, e := safe.Variable(run, text, aff); e.(type) {
		case nil:
			ret = val
		default:
			err = e
		case g.UnknownTarget, g.UnknownField:
			ret = nil // no such variable....
		}
	}
	return
}
