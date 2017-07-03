package rt

import (
	"io"
)

type ObjectFinder interface {
	FindObject(name string) (Object, bool)
}

type Runtime interface {
	// Writer writes standard output.
	io.Writer
	// PushWriter to set the active writer.
	PushWriter(io.Writer)
	// PopWriter to restore the writer active before the most recent PushWriter.
	PopWriter()

	Random(inclusiveMin, exclusiveMax int) int

	ObjectFinder
	PushScope(ObjectFinder)
	PopScope()

	// NewObject from the passed class. The object has no name and cannot be found via GetObject()
	NewObject(class string) (Object, error)
	// GetObject with the passed name.
	GetObject(name string) (Object, bool)
	// GetClass with the passed name.
	GetClass(name string) (Class, bool)
	// GetRelation with the passed name.
	GetRelation(name string) (Relation, bool)
}

func ExecuteList(run Runtime, x []Execute) (err error) {
	for _, s := range x {
		if e := s.Execute(run); e != nil {
			err = e
			break
		}
	}
	return
}
