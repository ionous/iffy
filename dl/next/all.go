package next

import "github.com/ionous/iffy/rt"

type AllTrue struct {
	Test []rt.BoolEval
}

func (a *AllTrue) GetBool(run rt.Runtime) (okay bool, err error) {
	i, cnt := 0, len(a.Test)
	for ; i < cnt; i++ {
		if ok, e := a.Test[i].GetBool(run); e != nil {
			err = e
			break
		} else if !ok {
			break
		}
	}
	okay = i == cnt
	return
}

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
