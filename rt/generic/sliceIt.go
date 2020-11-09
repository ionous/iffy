package generic

import (
	"math"
)

func SliceFloats(vs []float64) Iterator {
	return SliceIt(len(vs), func(i int) (ret Value, err error) {
		ret = FloatOf(vs[i])
		return
	})
}

func SliceStrings(vs []string) Iterator {
	return SliceIt(len(vs), func(i int) (ret Value, err error) {
		ret = StringOf(vs[i])
		return
	})
}

func SliceIt(size int, next func(int) (Value, error)) *sliceIt {
	return &sliceIt{size: size, next: next}
}

func ListIt(v Value) (ret *sliceIt) {
	if cnt, e := v.GetLen(); e != nil {
		ret = SliceIt(math.MaxInt64, func(i int) (Value, error) {
			return nil, e
		})
	} else {
		ret = SliceIt(cnt, func(i int) (Value, error) {
			return v.GetIndex(i)
		})
	}
	return
}

type sliceIt struct {
	at   int
	size int
	next func(int) (Value, error)
}

func (it *sliceIt) Remaining() int {
	return it.size - it.at
}

func (it *sliceIt) HasNext() bool {
	return it.at < it.size
}

func (it *sliceIt) GetNext() (ret Value, err error) {
	if !it.HasNext() {
		err = Overflow{it.at, it.size}
	} else if v, e := it.next(it.at); e != nil {
		err = e
	} else {
		ret = v
		it.at++
	}
	return
}
