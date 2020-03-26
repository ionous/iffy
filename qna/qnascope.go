package qna

// VariableScope access variables from the pool on the stack.
type VariableScope interface {
	// GetVariable writes the value at 'name' into the value pointed to by 'pv'.
	GetVariable(name string, pv interface{}) error
	// SetVariable writes the value of 'v' into the value at 'name'.
	SetVariable(name string, v interface{}) error
}

// LoopFactory helps add iterator parameters to current scope.
type LoopFactory struct {
	i       int
	upscope VariableScope
}

// NewLoop creates a factory for iterator variables while looping ( ex. through a series of objects. )
// index: counts the loop iterations, starting with 1.
// first: returns true for the first iteration, and false thereafter.
// last: returns true only during the last iteration.
func NewLoop(upscope VariableScope) *LoopFactory {
	return &LoopFactory{i: 0, upscope: upscope}
}

// NewScope creates a scope for this round of the loop;
// updates the internal counter for the next round of the loop.
func (l *LoopFactory) NextScope(hasNext bool) VariableScope {
	l.i++ // pre-inc, because while i starts at zero, the loop counter starts at one.
	return &loopScope{currIndex: l.i, hasNext: hasNext, upscope: l.upscope}
}

// UnknownLoopVariable error type for unknown variables while processing loops.
type UnknownLoopVariable string

// Error returns the name of the unknown variable.
func (e UnknownLoopVariable) Error() string { return string(e) }

// internal, implements VariableScope.
type loopScope struct {
	currIndex int
	hasNext   bool
	upscope   VariableScope
}

// GetVariable returns values for the iterator variables (index,first,last) and anything up-scope.
func (l *loopScope) GetVariable(n string, pv interface{}) (err error) {
	switch n {
	case "index":
		err = Assign(pv, l.currIndex)
	case "first":
		err = Assign(pv, l.currIndex == 1)
	case "last":
		err = Assign(pv, !l.hasNext)
	default:
		if l.upscope != nil {
			err = l.upscope.GetVariable(n, pv)
		} else {
			err = UnknownLoopVariable(n)
		}
	}
	return
}

// SetVariable transparently hands off to the upscope passed in NewLoop:
// loop variables are not writable.
func (l *loopScope) SetVariable(n string, v interface{}) (err error) {
	if l.upscope != nil {
		err = l.upscope.SetVariable(n, v)
	} else {
		err = UnknownLoopVariable(n)
	}
	return
}
