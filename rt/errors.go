package rt

import "github.com/ionous/errutil"

// error for GetVariable, SetVariable.
type UnknownVariable string

// Error returns the name of the unknown variable.
func (e UnknownVariable) Error() string {
	return errutil.Sprintf("variable not found %q", string(e))
}

// error for GetField, SetField
type UnknownField struct {
	Target, Field string
}

// Error returns the name of the unknown variable.
func (e UnknownField) Error() string {
	return errutil.Sprintf("field not found '%s.%s'", e.Target, e.Field)
}
