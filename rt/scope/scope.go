package scope

import "github.com/ionous/errutil"

// UnknownVariable error type for unknown variables while processing loops.
type UnknownVariable string

// Error returns the name of the unknown variable.
func (e UnknownVariable) Error() string {
	return errutil.Sprintf("unknown variable %q", string(e))
}

// EmptyScope allows use as a perpetually erroring scope.
type EmptyScope struct{}

func (EmptyScope) GetVariable(n string) (interface{}, error) {
	return nil, UnknownVariable(n)
}

func (EmptyScope) SetVariable(n string, v interface{}) error {
	return UnknownVariable(n)
}

type ReadOnly interface {
	GetVariable(string) (interface{}, error)
}
