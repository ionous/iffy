package core

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
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
	var idx int
	return stream.NewNumberStream(func() (ret interface{}, okay bool) {
		if idx < end {
			v := float64(idx)
			idx += step
			ret, okay = stream.Value(v)
		}
		return
	}), nil
}
