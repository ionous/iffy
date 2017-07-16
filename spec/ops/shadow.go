package ops

import (
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// ShadowClass implements a command for constructing pod-like types.
type ShadowClass struct {
	rtype r.Type
	evals map[string]r.Value
}

// GetObject for a shadow type generates an object from the evals specified.
// It is a constructor.
func (c *ShadowClass) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := run.NewObject(c.rtype.Name()); e != nil {
		err = e
	} else {
		// walk all the fields we recorded and pass them to the new object
		for k, eval := range c.evals {
			// Unpack evaluates an interface to get its resulting go value.
			if v, e := rt.Unpack(run, eval.Interface()); e != nil {
				err = e
				break
			} else if e := obj.SetValue(k, v); e != nil {
				err = e
				break
			}
		}
		if err == nil {
			ret = obj
		}
	}
	return
}

// Addr: all command interfaces are (normally) implemented as pointers.
// Spec carries around the target element, and has to take its address to make it into a pointer so that it matches the implementation.
// When using constructors, spec uses *ShadowClass as its target.
// We need to return the Value just of ourself.
func (c *ShadowClass) Addr() r.Value {
	return r.ValueOf(c)
}

// Type returns the type of the class for field walking.
// Compatible with reflect.Value
func (c *ShadowClass) Type() r.Type {
	return c.rtype
}

// Field returns the value of the requested field.
// Compatible with reflect.Value
// The spec will provide some type-safety on assignment to this value.
// FIX? one thing this cant handle is setting a state via an enumerated value.
// ex. TriState ( yes, no, maybe ) cmd.Param("yes").Value("true")
func (c *ShadowClass) Field(index int) (ret r.Value) {
	field := c.rtype.Field(index)
	if et, ok := evalFromType(field.Type); ok {
		val := r.New(r.TypeOf(et).Elem()).Elem()
		c.evals[field.Name] = val
		ret = val
	}
	return
}
