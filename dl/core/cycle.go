package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type CycleCounter struct {
	Name string `if:"id"`
	Curr int
}

// CycleText when called multiple times returns each of its inputs in turn
type CycleText struct {
	Id     string
	Values []rt.TextEval
}

func (l *CycleText) GetText(run rt.Runtime) (ret string, err error) {
	var curr int
	if obj, ok := run.GetObject(l.Id); !ok {
		err = errutil.New("couldnt find", l.Id)
	} else if e := obj.GetValue("curr", &curr); e != nil {
		err = e
	} else {
		next := (curr + 1) % len(l.Values)
		if e := obj.SetValue("curr", next); e != nil {
			err = e
		} else {
			at := l.Values[curr]
			ret, err = at.GetText(run)
		}
	}
	return
}
