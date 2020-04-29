package scope

import (
	"github.com/ionous/errutil"
)

// Parameters maps names to arbitrary evaluations.
type Parameters map[string]interface{}

// GetVariable returns the value at 'name', the caller is responsible for determining the type.
func (p Parameters) GetVariable(name string) (ret interface{}, err error) {
	if i, ok := p[name]; !ok {
		err = errutil.New("variable", name, "not found")
	} else {
		ret = i
	}
	return
}

// SetVariable writes the value of 'v' into the value at 'name'.
func (p Parameters) SetVariable(name string, v interface{}) (err error) {
	p[name] = v // FIX: any sort of validation? ex. ensure the value is baked ( ie. some sort of primitive or slice of primitives. )
	return
}
