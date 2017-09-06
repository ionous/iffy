package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"math"
)

type StoppingCounter struct {
	Name string `if:"id"`
	Curr float64
}

type StoppingText struct {
	Id     string
	Values []string
}

func (l *StoppingText) GetText(run rt.Runtime) (ret string, err error) {
	var curr int
	if obj, ok := run.FindObject(l.Id); !ok {
		err = errutil.New("couldnt find", l.Id)
	} else if e := run.GetValue(obj, "curr", &curr); e != nil {
		err = e
	} else {
		next := math.Min(float64(curr)+1, float64(len(l.Values)-1))
		if e := run.SetValue(obj, "curr", next); e != nil {
			err = e
		} else {
			ret = l.Values[curr]
		}
	}
	return
}
