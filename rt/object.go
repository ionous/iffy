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
	//
	GetValue(prop string, pv interface{}) error
	//
	SetValue(prop string, v interface{}) error
}
