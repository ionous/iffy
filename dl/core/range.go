package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
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

func (op *Range) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if vs, e := op.getNumList(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs
	}
	return
}

func (op *Range) getNumList(run rt.Runtime) (ret g.Value, err error) {
	if start, e := safe.GetOptionalNumber(run, op.Start, 1); e != nil {
		err = e
	} else if stop, e := safe.GetOptionalNumber(run, op.Stop, start.Float()); e != nil {
		err = e
	} else if step, e := safe.GetOptionalNumber(run, op.Step, 1); e != nil {
		err = e
	} else if step := step.Int(); step == 0 {
		err = errutil.New("Range error, step cannot be zero")
	} else {
		it := &rangeIt{start.Int(), stop.Int(), step}
		ret, err = g.CompactNumbers(it, nil)
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
func (it *rangeIt) GetNext() (ret g.Value, err error) {
	if !it.HasNext() {
		err = g.StreamExceeded
	} else {
		v := g.IntOf(it.next)
		it.next = it.next + it.step
		ret = v
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
