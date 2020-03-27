package qna

// VariableScope reads from and writes to a pool of named variables;
// the variables, their names, and initial values depend on the implementation and its context.
// Often, scopes are arranged in a stack with the newest scope checked for variables first, the oldest last.
type VariableScope interface {
	// GetVariable writes the value at 'name' into the value pointed to by 'pv'.
	GetVariable(name string, pv interface{}) error
	// SetVariable writes the value of 'v' into the value at 'name'.
	SetVariable(name string, v interface{}) error
}

// UnknownVariable error type for unknown variables while processing loops.
type UnknownVariable string

// Error returns the name of the unknown variable.
func (e UnknownVariable) Error() string { return string(e) }
