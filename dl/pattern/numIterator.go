package pattern

import (
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// implements chain.StreamIterator for multiple streams of numbers
type xnumIterator struct {
	run   rt.Runtime
	pat   *xNumListPattern
	order []int
	curr  int
}

func (k *xnumIterator) HasNextStream() bool {
	return k.curr < len(k.order)
}

func (k *xnumIterator) GetNextStream() (ret g.Iterator, err error) {
	if !k.HasNextStream() {
		err = g.StreamExceeded
	} else {
		ind := k.order[k.curr]
		if vs, e := safe.GetNumList(k.run, k.pat.Rules[ind]); e != nil {
			err = e
		} else {
			ret = g.ListIt(vs)
			k.curr++
		}
	}
	return
}
