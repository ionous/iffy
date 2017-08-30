package unique

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	r "reflect"
)

// Types implements a simple TypeRegistry.
type Types map[ident.Id]r.Type

// FindType by name.
func (reg Types) FindType(name string) (r.Type, bool) {
	id := ident.IdOf(name)
	rtype, ok := reg[id]
	return rtype, ok
}

// RegisterTypes implements TypeRegistry for a simple map.
func (reg Types) RegisterType(rtype r.Type) (err error) {
	id := ident.IdOf(rtype.Name())
	if was, ok := reg[id]; !ok {
		reg[id] = rtype
	} else if was != rtype {
		err = errutil.New("duplicate type", id, was, rtype)
	}
	return
}
