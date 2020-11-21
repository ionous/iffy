package generic

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

// panics if the value isnt a list
func ListIt(v Value) (ret *sliceIt) {
	return SliceIt(v.Len(), func(i int) (ret Value, _ error) {
		ret = v.Index(i)
		return
	})
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
