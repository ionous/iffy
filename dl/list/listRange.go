package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// Range generates a series of integers r[i] = (start + step*i) where i>=0
// Start and step default to 1, stop defaults to start;
// the inputs are truncated to produce whole numbers;
// a zero step returns an error.
//
// A positive step ends the series when the returned value would exceed stop.
// while a negative step ends before generating a value less than stop.
type Range struct {
	To     rt.NumberEval `if:"selector"`
	From   rt.NumberEval `if:"optional"`
	ByStep rt.NumberEval `if:"optional"`
}

func (*Range) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Function},
		Group:  "flow",
		Desc:   "Range of numbers: Generates a series of numbers.",
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
	if start, e := safe.GetOptionalNumber(run, op.From, 1); e != nil {
		err = e
	} else if stop, e := safe.GetOptionalNumber(run, op.To, start.Float()); e != nil {
		err = e
	} else if step, e := safe.GetOptionalNumber(run, op.ByStep, 1); e != nil {
		err = e
	} else if step := step.Int(); step == 0 {
		err = errutil.New("Range error, step cannot be zero")
	} else {
		ret = &ranger{start: start.Int(), stop: stop.Int(), step: step}
	}
	return
}

// ranger is a PanicValue where every method panics except type and affinity.
type ranger struct {
	g.Nothing
	start, stop, step int
}

// Affinity of a range is a number list.
func (n ranger) Affinity() affine.Affinity {
	return affine.NumList
}

// Type returns "range".
func (n ranger) Type() string {
	return "range"
}

// Index computes the i(th) step of the range.
func (n ranger) Index(i int) g.Value {
	v := n.start + i*n.step
	return g.IntOf(v)
}

// Len returns the total number of steps.
func (n ranger) Len() (ret int) {
	if diff := (n.stop - n.start + n.step); (n.step < 0) == (diff < 0) {
		ret = diff / n.step
	}
	return
}
