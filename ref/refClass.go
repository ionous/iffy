package ref

import (
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// note: for speed of some of these operations, we could cache the results in a pointer to a struct pooled in a map of rtype->RefClass.
// equality would work because the pointers are the same,
// just as equaity works by aliasing r.Type directly.
// note: r.Type is an interface, and go doesn't allow you to define methods for an interface type, even if you type define it as something else. ( which, frankly is just odd. you can add methods to typed primitives for goodness sake. why would typed interfaces be any different? they shouldn't be. )
type RefClass struct {
	r.Type
}

func MakeClass(rtype r.Type) RefClass {
	return RefClass{rtype}
}

// GetId returns the unique identifier for this classes.
func (c RefClass) GetId() string {
	return class.Id(c.Type)
}

// GetName returns a friendly name: spaces, no caps.
func (c RefClass) GetName() string {
	return class.FriendlyName(c.Type)
}

// GetParentType returns false for classes if no parent;
func (c RefClass) GetParent() (ret rt.Class, okay bool) {
	if cls, ok := class.Parent(c.Type); ok {
		ret, okay = MakeClass(cls), true
	}
	return
}

func (c RefClass) IsCompatible(name string) (okay bool) {
	return class.IsCompatible(c.Type, name)
}
