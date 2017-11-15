package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// ShadowClass provides a target which can construct a pod-like type.
// Each shadow instance can be assigned evals capable of filling the pod's corresponding fields; GetObject constructs the pod-type based on the assigned evals.
type ShadowClass struct {
	rtype  r.Type                 // the pod-type this instance constructs
	fields []FieldIndex           // direct access to all of the pod's fields; could be pooled/shared.
	slots  map[string]_ShadowSlot // evals used to construct a pod.
}

func Shadow(rtype r.Type) *ShadowClass {
	var fields []FieldIndex
	flatten(rtype, nil, &fields)
	return &ShadowClass{rtype, fields, make(map[string]_ShadowSlot)}
}

// GetObject constructs a pod-like object by evaluating the shadow's previously assigned evals.
// Each call generates a brand new object by reevaluating the shadow.
func (c *ShadowClass) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	obj := run.Emplace(r.New(c.rtype).Interface())
	// walk all the fields we recorded and pass them to the new object
	for k, slot := range c.slots {
		// evaluate and retrieve the resulting value
		if v, e := slot.unpack(run); e != nil {
			err = errutil.New("shadow class", c.rtype, "couldn't evaluate", k, e)
			break
		} else if e := obj.SetValue(k, v); e != nil {
			err = errutil.New("shadow class", c.rtype, "couldn't set value", k, e)
			break
		}
	}
	if err == nil {
		ret = obj
	}
	return
}

// Addr return the r.Value of ourself.
// All command interfaces are (normally) implemented as pointers.
// Spec carries around the target element, and has to take its address to make it into a pointer so that it matches the implementation.
// When using constructors, spec uses *ShadowClass as its target.
func (c *ShadowClass) Addr() r.Value {
	return r.ValueOf(c)
}

// Type, as per reflect.Value, returns the underlying pod-type for walking its fields.
func (c *ShadowClass) Type() r.Type {
	return c.rtype
}

// NumField, as per reflect.Value, returns the number of fields in the pod (and its parents).
func (c *ShadowClass) NumField() int {
	return len(c.fields)
}

// Field, as per reflect.Value, returns a slot corresponding to the indicated field.
// see also: FieldByIndex.
func (c *ShadowClass) Field(n int) (ret r.Value) {
	if n < len(c.fields) {
		ret = c.FieldByIndex(c.fields[n])
	}
	return
}

// FieldByName, as per reflect.Value, returns a slot corresponding to the named field in the shadow's pod type. see also: FieldByIndex.
func (c *ShadowClass) FieldByName(n string) (ret r.Value) {
	k := ident.IdOf(n)
	unique.WalkProperties(c.rtype, func(f *r.StructField, idx []int) (done bool) {
		if k == ident.IdOf(f.Name) {
			ret = c.getField(f, true)
			done = true
		}
		return
	})
	return
}

// FieldByIndex, as per reflect.Value, returns the slot of the indicated field.
// The slot allows the caller to poke in an eval which will be used to fill out the value of a pod during GetObject.
// As a side-effect it caches the slot; this as opposed to creating all fields in Shadow.
// FIX? one thing this cant handle is setting a state via an enumerated value.
// ex. TriState ( yes, no, maybe ) cmd.Param("yes").Value("true")
func (c *ShadowClass) FieldByIndex(n []int) r.Value {
	field := c.rtype.FieldByIndex(n)
	return c.getField(&field, true)
}

func (c *ShadowClass) getField(field *r.StructField, cache bool) (ret r.Value) {
	// determine what kind of eval can produce the passed type.
	if rtype := kindOf.EvalType(field.Type); rtype != nil {
		if x, ok := c.slots[field.Name]; ok {
			ret = x.rvalue
		} else if cache {
			// create an empty eval for the user to poke into
			rvalue := r.New(rtype).Elem()
			c.slots[field.Name] = _ShadowSlot{rtype, rvalue}
			ret = rvalue
		}
	}
	return
}
