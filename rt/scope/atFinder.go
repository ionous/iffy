package scope

import (
	"github.com/ionous/iffy/rt"
)

// Scope matches an single object using the name @.
type _AtFinder struct {
	obj rt.Object
}

func AtFinder(obj rt.Object) rt.ObjectFinder {
	return _AtFinder{obj}
}

// FindObject implements ObjectFinder.
func (l _AtFinder) FindObject(name string) (ret rt.Object, okay bool) {
	if name == "@" {
		ret, okay = l.obj, true
	}
	return
}
