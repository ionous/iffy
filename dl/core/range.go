package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

// Range generates a series of integers r[i] = (start + step*i) where i>=0
// Start defaults to 0, stop defaults to start, and step defaults to 1;
// the inputs are truncated to produce whole numbers;
// a zero step returns an error.
//
// A positive step ends the series when the returned value would exceed stop.
// while a negative step ends before generating a value less than stop.
type Range struct {
	Start, Stop, Step rt.NumberEval
}

func (*Range) Compose() composer.Spec {
	return composer.Spec{
		Name:  "range_over",
		Group: "flow",
		Desc:  "Range of numbers: Generates a series of numbers.",
	}
}

func (op *Range) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if start, e := rt.GetOptionalNumber(run, op.Start, 1); e != nil {
		err = e
	} else if stop, e := rt.GetOptionalNumber(run, op.Stop, start); e != nil {
		err = e
	} else if step, e := rt.GetOptionalNumber(run, op.Step, 1); e != nil {
		err = e
	} else if step == 0 {
		err = errutil.New("Range error, step cannot be zero")
	} else {
		it := &rangeIt{int(start), int(stop), int(step)}
		ret, err = rt.CompactNumbers(it, nil)
	}
	return
}

type rangeIt struct {
	next, stop, step int
}

// HasNext returns true if the iterator can be safely advanced.
func (it *rangeIt) HasNext() (okay bool) {
	return (it.step < 0 && it.next >= it.stop) ||
		(it.step > 0 && it.next <= it.stop)
}

// GetNumber advances the iterator.
func (it *rangeIt) GetNext() (ret rt.Value, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ret = generic.NewInt(it.next)
		it.next = it.next + it.step
	}
	return
}

// Remaining returns the remaining length of the stream.
func (it *rangeIt) Remaining() (ret int) {
	if cnt := 1 + (it.stop-it.next)/it.step; cnt > 0 {
		ret = cnt
	}
	return
}
