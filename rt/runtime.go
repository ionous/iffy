package rt

import (
	"io"
)

// ObjectFinder searches for objects by context.
type ObjectFinder interface {
	FindObject(name string) (Object, bool)
}

// Model describes the predefined world
type Model interface {
	// GetObject with the passed name.
	GetObject(name string) (Object, bool)
	// GetClass with the passed name.
	GetClass(name string) (Class, bool)
	// GetRelation with the passed name.
	GetRelation(name string) (Relation, bool)
}

type Runtime interface {
	// Model describes the predefined world
	Model
	// Writer writes standard output.
	io.Writer
	// PushWriter to set the active writer.
	PushWriter(io.Writer)
	// PopWriter to restore the writer active before the most recent PushWriter.
	PopWriter()
	// Random picks a pseudo-random value from a range. Can return any number including the lower bound, and up-to, but not including, the upper bound.
	Random(inclusiveMin, exclusiveMax int) int

	ObjectFinder
	PushScope(ObjectFinder)
	PopScope()

	// NewObject from the passed class. The object is anonymous. It has no name and cannot be found via GetObject()
	NewObject(class string) (Object, error)
	// Emplace adds an anonymous object to the runtime. The object has no name and cannot be found via GetObject().
	Emplace(mem interface{}) (Object, error)
	// GetAncestors returns a stream of objects starting with the passed object, then walking up whatever hierarchy the particular runtime implementation has defined.
	// E.g. parent-child containment.
	GetAncestors(Object) (ObjectStream, error)
}
