package scope

import (
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
)

// Scope matches an single object using the name @.
type _AtFinder struct {
	obj ref.Object
}

func AtFinder(obj ref.Object) rt.ObjectFinder {
	return _AtFinder{obj}
}

// FindObject implements ObjectFinder.
func (l _AtFinder) FindObject(name string) (ret ref.Object, okay bool) {
	if name == "@" {
		ret, okay = l.obj, true
	}
	return
}
