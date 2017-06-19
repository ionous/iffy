package core

import "github.com/ionous/iffy/rt"

type AllTrue struct {
	Test []rt.BoolEval
}

func (a *AllTrue) GetBool(run rt.Runtime) (ret bool, err error) {
	prelim := true
	for _, b := range a.Test {
		if ok, e := b.GetBool(run); e != nil {
			err = e
			break // see also any.
		} else if !ok {
			prelim = false
			break
		}
	}
	ret = prelim
	return
}

type AnyTrue struct {
	Test []rt.BoolEval
}

func (a *AnyTrue) GetBool(run rt.Runtime) (ret bool, err error) {
	prelim := false
	for _, b := range a.Test {
		if ok, e := b.GetBool(run); e != nil {
			err = e
			break
		} else if ok {
			prelim = true
			break
		}
	}
	ret = prelim
	return
}
