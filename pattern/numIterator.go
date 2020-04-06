package pattern

import (
	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// numIterator
type numIterator struct {
	run   rt.Runtime
	rules NumListRules
	order []int
	curr  int
}

func (k *numIterator) HasNext() bool {
	return k.curr < len(k.order)
}

func (k *numIterator) GetNext(pv interface{}) (err error) {
	if !k.HasNext() {
		err = stream.Exceeded
	} else if pit, ok := pv.(*rt.Iterator); !ok {
		err = assign.Mismatch("GetNext", pit, pv)
	} else {
		ind := k.order[k.curr]
		if it, e := rt.GetNumberStream(k.run, k.rules[ind]); e != nil {
			err = e
		} else {
			*pit = it
			k.curr++
		}
	}
	return
}
