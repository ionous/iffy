package rt

import (
	"github.com/ionous/iffy/ident"
	r "reflect"
)

// Class is an alias for reflect.Type
type Class r.Type

// Object represents a tangible or intangible piece of the game world.
type Object interface {
	// Id provides a somewhat unique identifier.
	Id() ident.Id
	// Type reflects on the variety of the object.
	Type() r.Type
	// GetValue returns the value of the named property via the passed pointer.
	// Errors if the property doesn't exist or if the property value cant be coerced to requested type.
	GetValue(prop string, pv interface{}) error
	// SetValue changes the value of the named property.
	// Errors if the property doesn't exist or if the incoming value
	// cant be coerced into a value compatible with the property's type.
	SetValue(prop string, v interface{}) error
}
