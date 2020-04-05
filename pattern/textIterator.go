package pattern

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// textIterator
type textIterator struct {
	run   rt.Runtime
	rules TextListRules
	order []int
	curr  int
}

func (k *textIterator) HasNext() bool {
	return k.curr < len(k.order)
}

func (k *textIterator) GetNext(pv interface{}) (err error) {
	if !k.HasNext() {
		err = stream.Exceeded
	} else if pit, ok := pv.(*rt.Iterator); !ok {
		err = errutil.Fmt("expected *iterator, got %T", pv)
	} else {
		ind := k.order[k.curr]
		if it, e := rt.GetTextStream(k.run, k.rules[ind]); e != nil {
			err = e
		} else {
			*pit = it
			k.curr++
		}
	}
	return
}
