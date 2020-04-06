package scope

import "github.com/ionous/iffy/rt"

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
func (l *loopScope) GetVariable(n string) (ret interface{}, err error) {
	switch n {
	case "index":
		ret = float64(l.currIndex)
	case "first":
		ret = bool(l.currIndex == 1)
	case "last":
		ret = bool(!l.hasNext)
	default:
		ret, err = l.val.GetVariable(n)
	}
	return
}

// SetVariable always returns an UnknownVariable error; iterator variables are not writable.
func (l *loopScope) SetVariable(n string, v interface{}) (err error) {
	return UnknownVariable(n)
}
