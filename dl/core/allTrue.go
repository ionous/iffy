package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// Always returns true
type Always struct{}

// AnyTrue returns true only when one of its specified tests returns true.
// It does not necessarily run all of the tests, it exits as soon as any test return true.
// An empty list returns false.
type AnyTrue struct {
	Test []rt.BoolEval
}

// AllTrue returns false only when one of its specified tests returns false.
// It does not necessarily run all of the tests, it exits as soon as any test return false.
// An empty list returns true.
type AllTrue struct {
	Test []rt.BoolEval
}

func (*Always) Compose() composer.Spec {
	return composer.Spec{
		Name:  "always",
		Group: "logic",
		Desc:  "Always: returns true always.",
	}
}

func (op *Always) GetBool(run rt.Runtime) (ret g.Value, err error) {
	ret = g.BoolOf(true)
	return
}

func (*AllTrue) Compose() composer.Spec {
	return composer.Spec{
		Name:  "all_true",
		Group: "logic",
		Spec:  "all true test: {+test|comma-and}",
		Desc:  "All True: returns true if all of the evaluations are true.",
	}
}

func (op *AllTrue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	// stop on the first statement to return false.
	if i, cnt, e := resolve(run, op.Test, false); e != nil {
		err = cmdError(op, e)
	} else if i < cnt {
		ret = g.False
	} else {
		ret = g.True // return true, resolve never found a false statement
	}
	return
}

func (*AnyTrue) Compose() composer.Spec {
	return composer.Spec{
		Name:  "any_true",
		Group: "logic",
		Desc:  "Any True: returns true if any of the evaluations are true.",
	}
}

func (op *AnyTrue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	// stop on the first statement to return true.
	if i, cnt, e := resolve(run, op.Test, true); e != nil {
		err = cmdError(op, e)
	} else if i < cnt {
		ret = g.True
	} else {
		ret = g.False // return false, resolve never found a true statement
	}
	return
}

func resolve(run rt.Runtime, evals []rt.BoolEval, breakOn bool) (i, cnt int, err error) {
	for i, cnt = 0, len(evals); i < cnt; i++ {
		if ok, e := safe.GetBool(run, evals[i]); e != nil {
			err = e
			break
		} else if ok.Bool() == breakOn {
			break
		}
	}
	return i, cnt, err
}
