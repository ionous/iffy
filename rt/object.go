package rt

import (
	"github.com/ionous/iffy/ident"
	r "reflect"
)

type Class r.Type

// Object represents a tangible or intangible piece of the game world.
type Object interface {
	// GetId returns a somewhat unique identifier.
	GetId() ident.Id
	// GetClass returns the variety of object.
	GetClass() Class
	// GetValue stores the value into the pointer pv.
	// Values include meta.Objects for relations and pointers, numbers, and text. For numbers, pv can be any numberic type: float64, int, etc.
	GetValue(name string, pv interface{}) error
	// GetValue can return error when the value violates a property constraint,
	// if the value is not of the requested type, or if the targeted property holder is read-only. Read-only values include the "many" side of a relation.
	SetValue(name string, v interface{}) error
}
