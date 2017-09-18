package rtm

import (
	"github.com/ionous/iffy/rt"
)

// ObjectScope provides name resolution by implementing rt.ObjectScope
type ObjectScope struct {
	topObject rt.Object
}

// TopObject returns the current top object, if any.
// By default. there is no top object.
func (os ObjectScope) TopObject() (ret rt.Object, okay bool) {
	if os.topObject != nil {
		ret, okay = os.topObject, true
	}
	return
}

// SetTop changes the current top object to the passed object.
func (os *ObjectScope) SetTop(top rt.Object) rt.Object {
	prev := os.topObject
	os.topObject = top
	return prev
}
