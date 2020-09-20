package scope

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

// LoopFactory or iterator variables while looping ( ex. through a series of objects. )
// index: counts the loop iterations, starting with 1.
// first: returns true for the first iteration, and false thereafter.
// last: returns true only during the last iteration.
type LoopFactory struct {
	varName string
	i       int
}

func LoopOver(run rt.Runtime, varName string, it rt.Iterator, do, other rt.Execute) (err error) {
	if hasNext := it.HasNext(); !hasNext {
		if e := rt.RunOne(run, other); e != nil {
			err = e
		}
	} else {
		lf := LoopFactory{varName: varName}
		for hasNext {
			if val, e := it.GetNext(); e != nil {
				err = e
			} else {
				hasNext = it.HasNext()
				// brings the names of an object's properties into scope for the duration of fn.
				run.PushScope(lf.NextScope(val, hasNext))
				e := rt.RunOne(run, do)
				run.PopScope()
				if e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

// NewScope creates a scope for this round of the loop;
// updates the internal counter for the next round of the loop.
func (l *LoopFactory) NextScope(varValue rt.Value, hasNext bool) rt.Scope {
	l.i++ // pre-inc, because while i starts at zero, the loop counter starts at one.
	return &loopScope{varName: l.varName, varValue: varValue, currIndex: l.i, hasNext: hasNext}
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
		ret = &generic.Int{Value: l.currIndex}
	case "first":
		ret = &generic.Bool{Value: l.currIndex == 1}
	case "last":
		ret = &generic.Bool{Value: !l.hasNext}
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
