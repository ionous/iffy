package scope

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

// EmptyScope allows use as a perpetually erroring scope.
type EmptyScope struct{}

func (EmptyScope) GetVariable(n string, pv interface{}) error {
	return UnknownVariable(n)
}

func (EmptyScope) SetVariable(n string, pv interface{}) error {
	return UnknownVariable(n)
}

type ReadOnly interface {
	GetVariable(string, interface{}) error
}
