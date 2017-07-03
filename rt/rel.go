package rt

import (
	"github.com/ionous/iffy/index"
)

// Relative describes the connection between objects.
type Relative interface {
	// GetValue mirrors Object.GetValue
	GetValue(name string, pv interface{}) error
	// SetValue mirrors Object.SetValue.
	// Changing the value of a relative property changes the status of the connection.
	SetValue(name string, v interface{}) error
}

// Relation connects Objects to each other in various ways.
type Relation interface {
	// GetId returns the unique identifier for this types.
	GetId() string
	// GetType of the relation: one-to-one to many-to-many.
	GetType() index.Type
	// Relate defines a connection between two objects.
	Relate(Object, Object) (Relative, error)
}
