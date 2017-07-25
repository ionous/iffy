package rt

import (
	"io"
)

// ExecuteList runs a block of statements.
type ExecuteList []Execute

func (x ExecuteList) Execute(run Runtime) (err error) {
	for _, s := range x {
		if e := s.Execute(run); e != nil {
			err = e
			break
		}
	}
	return
}

// Values for SetValues.
type Values map[string]interface{}

// SetValues to the passed object.
// FIX? add an optional map parameter to NewObject?
func SetValues(obj Object, values Values) (err error) {
	for name, v := range values {
		if e := obj.SetValue(name, v); e != nil {
			err = e
			break
		}
	}
	return
}

// Writer returns a new runtime that uses the passed writer for its output.
func Writer(run Runtime, w io.Writer) Runtime {
	return _Writer{run, w}
}

type _Writer struct {
	Runtime
	Writer io.Writer
}

func (l _Writer) Write(p []byte) (int, error) {
	return l.Writer.Write(p)
}
