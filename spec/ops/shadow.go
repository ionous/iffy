package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// ShadowClass implements a command for constructing pod-like types.
type ShadowClass struct {
	rtype r.Type
	slots map[string]_ShadowSlot
}

// note: nothing in the slot itself guarantees that the type and value are compatible.
// that's left up to spec/ops.
type _ShadowSlot struct {
	rtype  r.Type  // type of the slot
	rvalue r.Value // spec will .Set to this value
}

// GetObject for a shadow type generates an object from the slots specified.
// It is a constructor.
func (c *ShadowClass) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := run.NewObject(c.rtype.Name()); e != nil {
		err = errutil.New("shadow class", c.rtype, "couldn't create object")
	} else {
		// walk all the fields we recorded and pass them to the new object
		for k, slot := range c.slots {
			// Unpack evaluates an interface to get its resulting go value.
			if v, e := slot.unpack(run); e != nil {
				err = errutil.New("shadow class", c.rtype, "couldn't unpack", k, e)
				break
			} else if e := obj.SetValue(k, v); e != nil {
				err = errutil.New("shadow class", c.rtype, "couldn't set value", k, e)
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
	if rtype, ok := evalFromType(field.Type); ok {
		rvalue := r.New(rtype).Elem() // create an empty eval for the user to poke into
		c.slots[field.Name] = _ShadowSlot{rtype, rvalue}
		ret = rvalue
	}
	return
}
