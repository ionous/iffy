package generic

import "github.com/ionous/iffy/rt"

func SliceIt(size int, next func(int) rt.Value) *sliceIt {
	return &sliceIt{size: size, next: next}
}

func SliceFloats(vs []float64) *sliceIt {
	return SliceIt(len(vs), func(i int) rt.Value {
		return &Float{Value: vs[i]}
	})
}

func SliceStrings(vs []string) *sliceIt {
	return SliceIt(len(vs), func(i int) rt.Value {
		return &String{Value: vs[i]}
	})
}

type sliceIt struct {
	at   int
	size int
	next func(int) rt.Value
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
	} else {
		ret = it.next(it.at)
		it.at++
	}
	return
}
