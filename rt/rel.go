package rt

import (
	"github.com/ionous/iffy/index"
)

// Relation connects Objects to each other in various ways.
type Relation interface {
	// GetId returns the unique identifier for this types.
	GetId() string
	// GetType of the relation: one-to-one to many-to-many.
	GetType() index.Type
	// Relate defines a connection between two objects and a piece of data.
	// Returns previous data, if any.
	Relate(Object, Object, index.OnInsert) (bool, error)
	// Returns existing data, if any.
	GetRelative(Object, Object) (interface{}, bool)
	// hrm.
	GetTable() *index.Table
}
