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

// an object iterator that always fails
type EmptyObjects struct{}

func (EmptyObjects) HasNext() bool {
	return false
}

func (EmptyObjects) GetNext() (Object, error) {
	return nil, errutil.New("empty objects never has objects")
}
