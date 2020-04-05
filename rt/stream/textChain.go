package stream

import "github.com/ionous/iffy/rt"

type ChainText struct {
	evals rt.Iterator // iterate across streams
	items rt.Iterator // iterate across text
	err   error
}

// NewTextChain where evals.GetNext returns iterators whose GetNext returns text.
func NewTextChain(evals rt.Iterator) rt.Iterator {
	k := &ChainText{evals: evals, items: rt.EmptyStream(true)}
	k.items, k.err = k.advance()
	return k
}

func (k *ChainText) HasNext() bool {
	return k.err != Exceeded
}

func (k *ChainText) GetNext(pv interface{}) (err error) {
	if k.err != nil {
		err = k.err
	} else if e := k.items.GetNext(pv); e != nil {
		err = e
	} else {
		k.items, k.err = k.advance()
	}
	return
}

// can modify k.evals
func (k *ChainText) advance() (it rt.Iterator, err error) {
	for it = k.items; !it.HasNext(); {
		if e := k.evals.GetNext(&it); e != nil {
			err = e
			break
		}
	}
	return
}
