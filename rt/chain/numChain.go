package chain

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

type ChainNumbers struct {
	evals rt.Iterator // iterate across streams
	items rt.Iterator // iterate across numbers
	err   error
}

// NewNumberChain where evals.GetNext returns iterators whose GetNext returns numbers.
func NewNumberChain(evals rt.Iterator) rt.Iterator {
	k := &ChainNumbers{evals: evals, items: rt.EmptyStream(true)}
	k.items, k.err = k.advance()
	return k
}

func (k *ChainNumbers) HasNext() bool {
	return k.err != stream.Exceeded
}

func (k *ChainNumbers) GetNext(pv interface{}) (err error) {
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
func (k *ChainNumbers) advance() (it rt.Iterator, err error) {
	for it = k.items; !it.HasNext(); {
		if e := k.evals.GetNext(&it); e != nil {
			err = e
			break
		}
	}
	return
}
