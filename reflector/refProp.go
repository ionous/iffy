package reflector

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref"
	r "reflect"
)

type RefProp struct {
	id       string
	fieldIdx int //  index in parent rtype
	propType ref.PropertyType
}

// Property provides information on the fields of an object.
func (ref *RefProp) GetId() string {
	return ref.id
}

func (ref *RefProp) GetType() ref.PropertyType {
	return ref.propType
}

func (ref *RefProp) getFieldIndex() int {
	return ref.fieldIdx
}

// Categorize returns the property type for an go-type
func Categorize(rtype r.Type) (ret ref.PropertyType, err error) {
	switch k := rtype.Kind(); k {
	case r.Float64:
		ret = ref.Number
	case r.String:
		ret = ref.Text
	case r.Bool, r.Int:
		ret = ref.State
	case r.Ptr:
		// FIX: how far do you want to take this?
		ret = ref.Pointer
	case r.Slice:
		if elem, e := Categorize(rtype.Elem()); e != nil {
			err = e
		} else {
			ret = ref.Array | elem
		}
	default:
		err = errutil.New("unknown type", rtype, k.String())
	}
	return
}
