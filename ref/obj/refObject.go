package obj

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/enum"
	"github.com/ionous/iffy/ref/prop"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type RefObject struct {
	id    ident.Id // id of the object, blank if anonymous.
	value r.Value  // stores the concrete value. ex. Rock, not *Rock.
}

// Id returns the unique identifier for this Object.
// Blank for anonymous and temporary objects.
func (n RefObject) Id() ident.Id {
	return n.id
}

// Value holding the data for the object.
func (n RefObject) Value() r.Value {
	return n.value
}

// String representation of the object.
func (n RefObject) String() (ret string) {
	if n.id.IsValid() {
		ret = n.id.Name
	} else {
		ret = n.value.Type().Name()
	}
	return
}

// Type implements rt.Object.
func (n RefObject) Type() r.Type {
	return n.value.Type()
}

// Property implements rt.Object.
func (n RefObject) Property(name string) (ret rt.Property, okay bool) {
	rtype := n.Type()
	if path, idx := enum.PropertyPath(rtype, name); len(path) > 0 {
		pf := prop.MakeField(rtype.FieldByIndex(path), n.value.FieldByIndex(path))
		if idx < 0 {
			ret, okay = pf, true
		} else {
			ret, okay = prop.MakeState(pf, idx), true
		}
	}
	return
}
