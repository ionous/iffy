package core

import (
	"github.com/ionous/iffy/ref"
)

// Scope matches an single object using the name @.
type Scope struct {
	obj ref.Object
}

// NewScope creates a new object finder.
func NewScope(obj ref.Object) *Scope {
	return &Scope{obj}
}

// FindObject implements ObjectFinder.
func (l *Scope) FindObject(name string) (ret ref.Object, okay bool) {
	if name == "@" {
		ret, okay = l.obj, true
	}
	return
}
