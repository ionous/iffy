package scope

import (
	"github.com/ionous/errutil"
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
func (l *LoopFactory) NextScope(varName string, varValue rt.Value, hasNext bool) rt.Scope {
	l.i++ // pre-inc, because while i starts at zero, the loop counter starts at one.
	return &loopScope{varName: varName, varValue: varValue, currIndex: l.i, hasNext: hasNext}
}

// internal, implements Variable
type loopScope struct {
	varName   string
	varValue  rt.Value
	currIndex int
	hasNext   bool
}

// GetVariable returns values for the iterator variables (index,first,last) and anything up-
func (l *loopScope) GetVariable(n string) (ret rt.Value, err error) {
	switch n {
	case l.varName:
		ret = l.varValue
	case "index":
		ret = &rt.NumberValue{Value: float64(l.currIndex)}
	case "first":
		ret = &rt.BoolValue{Value: l.currIndex == 1}
	case "last":
		ret = &rt.BoolValue{Value: !l.hasNext}
	default:
		err = rt.UnknownVariable(n)
	}
	return
}

// SetVariable always returns an rt.UnknownVariable error; iterator variables are not writable.
func (l *loopScope) SetVariable(n string, v rt.Value) (err error) {
	switch n {
	case l.varName, "index", "first", "last":
		err = errutil.New("loop counters cant be changed")
	default:
		err = rt.UnknownVariable(n)
	}
	return
}
