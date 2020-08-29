package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
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

func (a *Always) GetBool(run rt.Runtime) (okay bool, err error) {
	return true, nil
}

func (*AllTrue) Compose() composer.Spec {
	return composer.Spec{
		Name:  "all_true",
		Group: "logic",
		Desc:  "All True: returns true if all of the evaluations are true.",
	}
}

func (a *AllTrue) GetBool(run rt.Runtime) (okay bool, err error) {
	eval := a.Test
	i, cnt := 0, len(eval)
	for ; i < cnt; i++ {
		if ok, e := GetBool(run, eval[i]); e != nil {
			err = e
			break
		} else if !ok {
			break
		}
	}
	if i == cnt {
		okay = true
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

func (a *AnyTrue) GetBool(run rt.Runtime) (okay bool, err error) {
	for _, b := range a.Test {
		if ok, e := rt.GetBool(run, b); e != nil {
			err = e
			break
		} else if ok {
			okay = true
			break
		}
	}
	return
}
