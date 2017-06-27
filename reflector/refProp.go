package reflector

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type RefProp struct {
	id       string
	fieldIdx int //  index in parent rtype
	propType rt.PropertyType
}

// Property provides information on the fields of an object.
func (p *RefProp) GetId() string {
	return p.id
}

func (p *RefProp) GetType() rt.PropertyType {
	return p.propType
}

func (p *RefProp) getFieldIndex() int {
	return p.fieldIdx
}

// Categorize returns the property type for an go-type
func Categorize(rtype r.Type) (ret rt.PropertyType, err error) {
	switch k := rtype.Kind(); k {
	case r.Float64:
		ret = rt.Number
	case r.String:
		ret = rt.Text
	case r.Bool, r.Int:
		ret = rt.State
	case r.Ptr:
		// FIX: how far do you want to take this?
		ret = rt.Pointer
	case r.Slice:
		if elem, e := Categorize(rtype.Elem()); e != nil {
			err = e
		} else {
			ret = rt.Array | elem
		}
	default:
		err = errutil.New("unknown type", rtype, k.String())
	}
	return
}
