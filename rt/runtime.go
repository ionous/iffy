package rt

import (
	"io"
)

// Pluralize turns single words into their plural variants.
type Pluralize interface {
	Pluralize(single string) string
}

// ObjectFinder searches for objects by context.
type ObjectFinder interface {
	FindObject(name string) (Object, bool)
}

// Model interacts with the predefined world.
type Model interface {
	// GetObject with the passed name.
	GetObject(name string) (Object, bool)
	// GetRelation with the passed name.
	GetRelation(name string) (Relation, bool)
	// Emplace wraps the passed value as an anonymous object.
	Emplace(i interface{}) Object
}

// Ancestors customizes the parent-child event hierarchy.
type Ancestors interface {
	// GetAncestors returns a stream of objects starting with the parent of the passed object, walking up whatever hierarchy the particular runtime implementation has defined.
	GetAncestors(Runtime, Object) (ObjectStream, error)
}

// Runtime environment for an in-progress game.
type Runtime interface {
	Model
	io.Writer
	Random(inclusiveMin, exclusiveMax int) int
	ObjectFinder
	Ancestors
	Patterns
	Pluralize
}
