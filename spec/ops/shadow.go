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
	}
	return
}

// Type returns the type of the class for field walking.
// Compatible with _SpecWriter
func (c *ShadowClass) Type() r.Type {
	return c.rtype
}

// Field returns the value of the request field
// Compatible with _SpecWriter
func (c *ShadowClass) Field(index int) (ret r.Value) {
	field := c.rtype.Field(index)
	// so, this is interesting:
	// we an eval which is capable of computing the pod-like type of the requested field
	// that's basically like specifying some literal for a command's eval.
	// we can re-use literally for that.
	zero := r.Zero(field.Type).Interface()
	if eval, ok := literally(field.Type, zero); ok {
		reveal := r.ValueOf(eval)
		c.evals[field.Name] = reveal
		ret = reveal
	}
	return
}
