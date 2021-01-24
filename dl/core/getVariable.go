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
	Name string `if:"selector"`
}

// Compose implements composer.Composer
func (*Var) Compose() composer.Spec {
	return composer.Spec{
		Name:   "get_var",
		Group:  "variables",
		Desc:   "Get Variable: Return the value of the named variable.",
		Fluent: &composer.Fluid{Name: "var", Role: composer.Function},
	}
}

// Affinity helps implement Assignment, but always returns ""
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

func (op *Var) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.NumList)
}

func (op *Var) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.TextList)
}

func (op *Var) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.RecordList)
}

func (op *Var) getVar(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Name, aff); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
