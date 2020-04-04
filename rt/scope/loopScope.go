package scope

import (
	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/rt"
)

// LoopFactory or iterator variables while looping ( ex. through a series of objects. )
// index: counts the loop iterations, starting with 1.
// first: returns true for the first iteration, and false thereafter.
// last: returns true only during the last iteration.
type LoopFactory struct {
	i int
}

// NewScope creates a scope for this round of the loop;
// updates the internal counter for the next round of the loop.
func (l *LoopFactory) NextScope(val ReadOnly, hasNext bool) rt.VariableScope {
	l.i++ // pre-inc, because while i starts at zero, the loop counter starts at one.
	return &loopScope{val: val, currIndex: l.i, hasNext: hasNext}
}

// internal, implements Variable
type loopScope struct {
	val       ReadOnly
	currIndex int
	hasNext   bool
}

// GetVariable returns values for the iterator variables (index,first,last) and anything up-
func (l *loopScope) GetVariable(n string, pv interface{}) (err error) {
	switch n {
	case "index":
		err = assign.ToFloat(pv, float64(l.currIndex))
	case "first":
		err = assign.ToBool(pv, l.currIndex == 1)
	case "last":
		err = assign.ToBool(pv, !l.hasNext)
	default:
		err = l.val.GetVariable(n, pv)
	}
	return
}

// SetVariable always returns an UnknownVariable error; iterator variables are not writable.
func (l *loopScope) SetVariable(n string, v interface{}) (err error) {
	return UnknownVariable(n)
}
