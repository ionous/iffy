package rt

import (
	"io"
)

func ExecuteList(run Runtime, x []Execute) (err error) {
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

func Write(run Runtime, w io.Writer, cb func() error) error {
	run.PushWriter(w)
	e := cb()
	run.PopWriter()
	return e
}
