package pattern

import (
	"github.com/ionous/iffy/rt"
)

// implements chain.StreamIterator for multiple streams of text
type textIterator struct {
	run   rt.Runtime
	pat   *TextListPattern
	order []int
	curr  int
}

func (k *textIterator) HasNextStream() bool {
	return k.curr < len(k.order)
}

func (k *textIterator) GetNextStream() (ret rt.Iterator, err error) {
	if !k.HasNextStream() {
		err = rt.StreamExceeded
	} else {
		ind := k.order[k.curr]
		if it, e := rt.GetTextStream(k.run, k.pat.Rules[ind]); e != nil {
			err = e
		} else {
			ret = it
			k.curr++
		}
	}
	return
}
