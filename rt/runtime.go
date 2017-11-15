package rt

import (
	"io"
	r "reflect"
)

// Pluralize turns single words into their plural variants.
type Pluralize interface {
	Pluralize(single string) string
}

// ObjectScope defines a system for name resolution
// using a single object as a set of local names
// similar to the parameters of a function call.
// see also: ScopeBlock()
type ObjectScope interface {
	TopObject() (Object, bool)
	SetTop(Object) (prev Object)
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

type Output interface {
	Writer() io.Writer
	SetWriter(io.Writer) (prev io.Writer)
}

// Runtime environment for an in-progress game.
type Runtime interface {
	Model
	Output
	Random(inclusiveMin, exclusiveMax int) int
	ObjectScope
	Ancestors
	Pattern
	Pluralize
	Pack(pdst, src r.Value) error
}

// WritersBlock applies a writer to the runtime for the duration of fn.
// If the writer also implements io.Closer and fn is free of errors,
// its Close() will be called, and the result of Close returned.
func WritersBlock(run Runtime, w io.Writer, fn func() error) (err error) {
	prev := run.SetWriter(w)
	e := fn()
	run.SetWriter(prev)
	if e != nil {
		err = e
	} else if closer, ok := w.(io.Closer); ok {
		err = closer.Close()
	}
	return
}

// ScopeBlock brings the names of an object's properties into scope for the duration of fn.
func ScopeBlock(run Runtime, top Object, fn func()) {
	prev := run.SetTop(top)
	fn()
	run.SetTop(prev)
}

func Determine(run Runtime, p interface{}) error {
	return run.ExecuteMatching(run.Emplace(p))
}
