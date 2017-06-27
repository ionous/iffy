package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type ShuffleText struct {
	Id     string
	Values []string
}

type ShuffleCounter struct {
	Name    string `if:"id"`
	Curr    float64
	Indices []float64
}

func (l *ShuffleText) GetText(run rt.Runtime) (ret string, err error) {
	if obj, ok := run.FindObject(l.Id); !ok {
		err = errutil.New("couldnt find", l.Id)
	} else {
		var curr int
		var indices []float64
		if e := obj.GetValue("curr", &curr); e != nil {
			err = e
		} else if e := obj.GetValue("indices", &indices); e != nil {
			err = e
		} else if cnt := len(l.Values); cnt > 0 {
			if len(indices) == 0 {
				indices = make([]float64, cnt)
				for i := 0; i < cnt; i++ {
					indices[i] = float64(i)
				}
				if e := obj.SetValue("indices", indices); e != nil {
					err = e
				}
			}
			if err == nil {
				if curr >= cnt {
					curr = 0 // wrap for the next sequence
				}
				j := run.Random(curr, cnt)
				if curr != j { // switch if they are different locations.
					indices[curr], indices[j] = indices[j], indices[curr]
				}
				if e := obj.SetValue("curr", curr+1); e != nil {
					err = e
				} else {
					sel := int(indices[curr])
					ret = l.Values[sel]

				}
			}
		}
	}
	return
}
