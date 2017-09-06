package rt

import (
	"github.com/ionous/iffy/ident"
	r "reflect"
)

// Class is an alias for reflect.Type
type Class r.Type

// Object represents a tangible or intangible piece of the game world.
type Object interface {
	// Id returns a somewhat unique identifier.
	Id() ident.Id
	// Type returns the variety of object.
	Type() r.Type
	// Property returns the named property, if it exists.
	Property(name string) (Property, bool)
}

type Property interface {
	// Id returns a somewhat unique identifier.
	Id() ident.Id
	// Type returns the variety of property.
	Type() r.Type
	// GetValue stores the value of this property into the pointer pv.
	// Values include meta.Objects for relations and pointers, numbers, and text. For numbers, pv can be any numberic type: float64, int, etc.
	GetValue(pv interface{}) error
	// SetValue stores the passed value into this property.
	// This can return error when the value violates a property constraint,
	// or if the value is not of the requested type.
	SetValue(v interface{}) error
}
