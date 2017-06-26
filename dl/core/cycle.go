package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type CycleCounter struct {
	Name string `if:"id"`
	Curr float64
}

type CycleText struct {
	Id     string
	Values []string
}

func (l *CycleText) GetText(run rt.Runtime) (ret string, err error) {
	var curr int
	if obj, ok := run.FindObject(l.Id); !ok {
		err = errutil.New("couldnt find", l.Id)
	} else if e := obj.GetValue("curr", &curr); e != nil {
		err = e
	} else {
		next := (curr + 1) % len(l.Values)
		if e := obj.SetValue("curr", next); e != nil {
			err = e
		} else {
			ret = l.Values[curr]
		}
	}
	return
}
