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
	// GetClass with the passed name.
	GetClass(name string) (Class, bool)
	// GetRelation with the passed name.
	GetRelation(name string) (Relation, bool)
	// GetValue sets the value of the passed pointer to the value of the named property in the passed object.
	GetValue(obj Object, name string, pv interface{}) error
	// SetValue sets the named property in the passed object to the value.
	SetValue(obj Object, name string, v interface{}) error
}

// Ancestors customizes the parent-child event hierarchy.
type Ancestors interface {
	// GetAncestors returns a stream of objects starting with the parent of the passed object, walking up whatever hierarchy the particular runtime implementation has defined.
	GetAncestors(Runtime, Object) (ObjectStream, error)
}

type Runtime interface {
	// Model describes the predefined world
	Model
	// Writer writes standard output.
	io.Writer
	// Random picks a pseudo-random value from a range. Can return any number including the lower bound, and up-to, but not including, the upper bound.
	Random(inclusiveMin, exclusiveMax int) int

	ObjectFinder

	// Emplace adds an anonymous object to the runtime. The object cannot be found via GetObject().
	Emplace(mem interface{}) Object
	//
	Ancestors
	// Patterns for pattern matching, iffy's equivalent of user methods and functions.
	Patterns
	// Pluralize for pluralization of printed nouns.
	Pluralize
}
