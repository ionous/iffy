package chain

import (
	g "github.com/ionous/iffy/rt/generic"
)

// iterate across multiple streams
type StreamIterator interface {
	HasNextStream() bool
	GetNextStream() (g.Iterator, error) // return the first iterator of the next stream
}

type chain struct {
	streams StreamIterator // iterate across streams
	items   g.Iterator     // iterate across numbers
	err     error
}

// NewStreamOfStreams turns multiple streams into a single iterator
func NewStreamOfStreams(streams StreamIterator) g.Iterator {
	k := &chain{streams: streams, items: g.EmptyStream(true)}
	k.items, k.err = k.advance()
	return k
}

func (k *chain) HasNext() bool {
	return k.err != g.StreamExceeded
}

func (k *chain) GetNext() (ret g.Value, err error) {
	if k.err != nil {
		err = k.err
	} else if next, e := k.items.GetNext(); e != nil {
		err = e
	} else {
		k.items, k.err = k.advance()
		ret = next
	}
	return
}

// attempt to move the iterator forward,
// at the end of its stream, move to the next stream
// return the new iterator
// can advance k.streams,
func (k *chain) advance() (it g.Iterator, err error) {
	for it = k.items; !it.HasNext(); {
		if next, e := k.streams.GetNextStream(); e != nil {
			err = e
			break
		} else {
			it = next
			// continue so we can text HasNext
		}
	}
	return
}
