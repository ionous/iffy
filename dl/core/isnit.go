package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// IsTrue transparently returns a boolean eval.
// It exists to help smooth the use of command expressions:
// eg. "is" {some expression}
type IsTrue struct {
	Test rt.BoolEval
}

// IsNotTrue returns the opposite of a boolean eval.
type IsNotTrue struct {
	Test rt.BoolEval
}

func (*IsTrue) Compose() composer.Spec {
	return composer.Spec{
		Name:  "is_true",
		Group: "logic",
		Desc:  "Is True: Transparently returns the result of a boolean expression.",
		Spec:  "{test:bool_eval} is true",
	}
}

func (op *IsTrue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetBool(run, op.Test); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (*IsNotTrue) Compose() composer.Spec {
	return composer.Spec{
		Name:  "not",
		Group: "logic",
		Desc:  "Is Not: Returns the opposite value.",
	}
}

func (op *IsNotTrue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetBool(run, op.Test); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(!val.Bool())
	}
	return
}
