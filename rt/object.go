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

// Property represents, roughly, the value of a field in an object, but it can also represent the status of an object state, or the relationship between objects.
type Property interface {
	// Id returns a somewhat unique identifier.
	Id() ident.Id
	// Type of the slot ( related to, but not always the same type as the value. )
	Type() r.Type
	// Value of the property.
	Value() interface{}
	// SetValue to change the property.
	SetValue(interface{}) error
}
