package core

import (
	"github.com/ionous/iffy/rt"
)

// Range generates a series of float values.
// FIX: look more at python frange
// FIX: add tests
type Range struct {
	Start, End, Step float64
}

func (l *Range) GetNumberStream(rt.Runtime) (rt.NumberStream, error) {
	start := int(l.Start)
	end := int(l.End)
	step := int(l.Step)
	if step == 0 {
		step = 1
	}
	if end == 0 {
		end = start
		start = 0
	}
	return &RangeIt{
		idx:  start,
		end:  end,
		step: step,
	}, nil
}

type RangeIt struct {
	idx, end, step int
}

func (it *RangeIt) HasNext() bool {
	return it.idx < it.end
}

func (it *RangeIt) GetNext() (ret float64, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ret = float64(it.idx)
		it.idx += it.step
	}
	return
}
