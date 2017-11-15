package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type StoppingCounter struct {
	Name string `if:"id"`
	Curr int
}

// StoppingText options returns its values in-order, once per evaluation, until the options are exhausted, then the last one is kept.
// As a special case, if there is only one option: it gets returned once, followed by the empty string forever after.
type StoppingText struct {
	Id     string
	Values []rt.TextEval
}

func (l *StoppingText) GetText(run rt.Runtime) (ret string, err error) {
	var curr int
	if obj, ok := run.GetObject(l.Id); !ok {
		err = errutil.New("couldnt find", l.Id)
	} else if e := obj.GetValue("curr", &curr); e != nil {
		err = e
	} else {
		switch cnt := len(l.Values); cnt {
		case 0:
		case 1:
			// only get text if at first:
			if curr == 0 {
				ret, err = l.update(run, obj, curr)
			}
		default:
			if curr >= cnt {
				curr = curr - 1
			}
			ret, err = l.update(run, obj, curr)
		}
	}
	return
}
func (l *StoppingText) update(run rt.Runtime, obj rt.Object, curr int) (ret string, err error) {
	if e := obj.SetValue("curr", curr+1); e != nil {
		err = e
	} else {
		at := l.Values[curr]
		ret, err = at.GetText(run)
	}
	return
}
