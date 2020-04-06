package next

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
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

func (*LenOfNumbers) Compose() composer.Spec {
	return composer.Spec{
		Name:  "len",
		Group: "format",
		Desc:  "Length of Number List: Determines the number of elements in a list of numbers.",
	}
}

func (op *LenOfNumbers) GetNumber(run rt.Runtime) (ret float64, err error) {
	// FIX? maybe the evals themselves should implement Count and not the activated stream.
	if elems, e := rt.GetNumberStream(run, op.Elems); e != nil {
		err = e
	} else if l, ok := elems.(rt.StreamCount); !ok {
		err = assign.Mismatch("unknown number list", l, elems)
	} else {
		ret = float64(l.Remaining())
	}
	return
}

func (*LenOfTexts) Compose() composer.Spec {
	return composer.Spec{
		Name:  "len",
		Group: "format",
		Desc:  "Length of Text List: Determines the number of text elements in a list.",
	}
}

func (op *LenOfTexts) GetNumber(run rt.Runtime) (ret float64, err error) {
	// FIX? maybe the evals themselves should implement Count and not the activated stream.
	if elems, e := rt.GetTextStream(run, op.Elems); e != nil {
		err = e
	} else if l, ok := elems.(rt.StreamCount); !ok {
		err = assign.Mismatch("unknown text list", l, elems)
	} else {
		ret = float64(l.Remaining())
	}
	return
}

func (*Range) Compose() composer.Spec {
	return composer.Spec{
		Name:  "range_over",
		Group: "flow",
		Desc:  "Range of numbers: Generates a series of numbers.",
	}
}

func (op *Range) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
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
func (it *rangeIt) GetNext(pv interface{}) (err error) {
	if !it.HasNext() {
		err = stream.Exceeded
	} else {
		now, next := it.next, it.next+it.step
		if e := assign.FloatPtr(pv, float64(now)); e != nil {
			err = e
		} else {
			it.next = next
		}
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
