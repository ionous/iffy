package rt

import (
	"github.com/ionous/iffy/ref"
	"io"
)

type ObjectFinder interface {
	FindObject(name string) (ref.Object, bool)
}

type Runtime interface {
	ref.Model

	Random(inclusiveMin, exclusiveMax int) int

	ObjectFinder
	PushScope(ObjectFinder)
	PopScope()

	io.Writer
	PushWriter(io.Writer)
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
