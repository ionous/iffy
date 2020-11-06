package chain

import (
	"math"

	"github.com/ionous/iffy/rt"
)

func SliceIt(size int, next func(int) (rt.Value, error)) *sliceIt {
	return &sliceIt{size: size, next: next}
}

func ListIt(v rt.Value) (ret *sliceIt) {
	if cnt, e := v.GetLen(); e != nil {
		ret = SliceIt(math.MaxInt64, func(i int) (rt.Value, error) {
			return nil, e
		})
	} else {
		ret = SliceIt(cnt, func(i int) (rt.Value, error) {
			return v.GetIndex(i)
		})
	}
	return
}

type sliceIt struct {
	at   int
	size int
	next func(int) (rt.Value, error)
}

func (it *sliceIt) Remaining() int {
	return it.size - it.at
}

func (it *sliceIt) HasNext() bool {
	return it.at < it.size
}

func (it *sliceIt) GetNext() (ret rt.Value, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else if v, e := it.next(it.at); e != nil {
		err = e
	} else {
		ret = v
		it.at++
	}
	return
}
