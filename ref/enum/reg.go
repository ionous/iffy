package enum

import (
	"github.com/ionous/iffy/ref/unique"
	r "reflect"
)

// Register the passed table data into the global registry.
// Panics on error.
func Register(ptr interface{}, s string, i []uint8) {
	if _, e := registry.Register(ptr, s, i); e != nil {
		panic(e)
	}
}

// Enumerate the passed type, caching to the global registry.
func Enumerate(rtype r.Type) []string {
	cs, _ := registry.Enumerate(rtype)
	return cs
}

// Registry caches the choices generated from enumerated types.
type Registry map[r.Type][]string

var registry = make(Registry)

// Register creates an Enum by using table data provided by stringer.
func (reg Registry) Register(ptr interface{}, s string, i []uint8) (ret []string, err error) {
	if rtype, e := unique.TypePtr(r.Int, ptr); e != nil {
		err = e
	} else if cs, ok := reg[rtype]; ok {
		ret = cs
	} else if cs, e := Compact(rtype, s, i); e != nil {
		err = e
	} else {
		ret, reg[rtype] = cs, cs
	}
	return
}

// Enumerate creates an Enum by probing possible values of the passed type.
func (reg Registry) Enumerate(rtype r.Type) (ret []string, err error) {
	if cs, ok := reg[rtype]; ok {
		ret = cs
	} else if cs, e := Stringify(rtype); e != nil {
		err = e
	} else {
		ret, reg[rtype] = cs, cs
	}
	return
}
