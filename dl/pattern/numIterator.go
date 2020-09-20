package pattern

import (
	"github.com/ionous/iffy/rt"
)

// implements chain.StreamIterator for multiple streams of numbers
type numIterator struct {
	run   rt.Runtime
	pat   *NumListPattern
	order []int
	curr  int
}

func (k *numIterator) HasNextStream() bool {
	return k.curr < len(k.order)
}

func (k *numIterator) GetNextStream() (ret rt.Iterator, err error) {
	if !k.HasNextStream() {
		err = rt.StreamExceeded
	} else {
		ind := k.order[k.curr]
		if it, e := rt.GetNumberStream(k.run, k.pat.Rules[ind]); e != nil {
			err = e
		} else {
			ret = it
			k.curr++
		}
	}
	return
}
