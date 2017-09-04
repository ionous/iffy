package rt

import (
	"github.com/ionous/errutil"
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

// an object iterator that always fails
type EmptyObjects struct{}

func (EmptyObjects) HasNext() bool {
	return false
}

func (EmptyObjects) GetNext() (Object, error) {
	return nil, errutil.New("empty objects never has objects")
}
