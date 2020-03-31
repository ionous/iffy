package next

import "github.com/ionous/iffy/rt"

// AllTrue returns false only when one of its specified tests returns false.
// It does not necessarily run all of the tests, it exits as soon as any test return false.
// An empty list returns true.
type AllTrue struct {
	Test []rt.BoolEval
}

func (a *AllTrue) GetBool(run rt.Runtime) (okay bool, err error) {
	if anyFalse, e := a.anyFalse(run); e != nil {
		err = e
	} else {
		okay = !anyFalse
	}
	return
}

func (a *AllTrue) anyFalse(run rt.Runtime) (ret bool, err error) {
	for _, b := range a.Test {
		if ok, e := b.GetBool(run); e != nil {
			err = e
			break
		} else if !ok {
			ret = true
			break
		}
	}
	return
}

// AnyTrue returns true only when one of its specified tests returns true.
// It does not necessarily run all of the tests, it exits as soon as any test return true.
// An empty list returns false.
type AnyTrue struct {
	Test []rt.BoolEval
}

func (a *AnyTrue) GetBool(run rt.Runtime) (okay bool, err error) {
	for _, b := range a.Test {
		if ok, e := b.GetBool(run); e != nil {
			err = e
			break
		} else if ok {
			okay = true
			break
		}
	}
	return
}
