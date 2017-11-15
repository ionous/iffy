package obj

import (
	"github.com/ionous/iffy/ident"
	r "reflect"
)

// Property represents, roughly, the value of a field in an object, but it can also represent the status of an object state, or the relationship between objects.
type Property interface {
	// Id returns a somewhat unique identifier.
	Id() ident.Id
	// Type of the slot ( related to, but not always the same type as the value. )
	Type() r.Type
	// Value of the property.
	Value() r.Value
	// SetValue to change the property.
	SetValue(r.Value) error
}
