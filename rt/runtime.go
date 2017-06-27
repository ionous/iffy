package rt

import (
	"io"
)

type ObjectFinder interface {
	FindObject(name string) (Object, bool)
}

type Runtime interface {
	Random(inclusiveMin, exclusiveMax int) int

	ObjectFinder
	PushScope(ObjectFinder)
	PopScope()

	// NewObject from the passed class. The object has no name and cannot be found via GetObject()
	NewObject(class string) (Object, error)
	// GetObject returns the object with the passed name.
	GetObject(name string) (Object, bool)
	// GetClass returns the class with the passed name.
	GetClass(name string) (Class, bool)

	// Writer writes standard output.
	io.Writer
	// PushWriter to set the active writer.
	PushWriter(io.Writer)
	// PopWriter to restore the writer active before the most recent PushWriter.
	PopWriter()
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
