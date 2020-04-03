package next

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
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

type LenOfNumbers struct {
	Elems rt.NumListEval
}

type LenOfTexts struct {
	Elems rt.TextListEval
}

func (op *LenOfNumbers) GetNumber(run rt.Runtime) (ret float64, err error) {
	// FIX? maybe the evals themselves should implement Count and not the activated stream.
	if elems, e := rt.GetNumberStream(run, op.Elems); e != nil {
		err = e
	} else if l, ok := elems.(rt.StreamCount); !ok {
		err = errutil.Fmt("unknown number list %T", elems)
	} else {
		ret = float64(l.Count())
	}
	return
}

func (op *LenOfTexts) GetNumber(run rt.Runtime) (ret float64, err error) {
	// FIX? maybe the evals themselves should implement Count and not the activated stream.
	if elems, e := rt.GetTextStream(run, op.Elems); e != nil {
		err = e
	} else if l, ok := elems.(rt.StreamCount); !ok {
		err = errutil.Fmt("unknown text list %T", elems)
	} else {
		ret = float64(l.Count())
	}
	return
}

func (op *Range) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	if start, e := rt.GetOptionalNumber(run, op.Start, 1); e != nil {
		err = e
	} else if stop, e := rt.GetOptionalNumber(run, op.Stop, start); e != nil {
		err = e
	} else if step, e := rt.GetOptionalNumber(run, op.Step, 1); e != nil {
		err = e
	} else if step == 0 {
		err = errutil.New("Range error, step cannot be zero")
	} else {
		ret = &rangeIt{int(start), int(stop), int(step)}
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
func (it *rangeIt) GetNumber() (ret float64, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		now, next := it.next, it.next+it.step
		ret, it.next = float64(now), next
	}
	return
}

// Count returns the remaining length of the stream.
func (it *rangeIt) Count() (ret int) {
	if cnt := 1 + (it.stop-it.next)/it.step; cnt > 0 {
		ret = cnt
	}
	return
}
