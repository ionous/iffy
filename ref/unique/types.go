package unique

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	r "reflect"
)

// Types implements a simple TypeRegistry.
type Types map[string]r.Type

// FindType by name.
func (reg Types) FindType(name string) (r.Type, bool) {
	id := id.MakeId(name)
	rtype, ok := reg[id]
	return rtype, ok
}

// RegisterTypes implements TypeRegistry for a simple map.
func (reg Types) RegisterType(rtype r.Type) (err error) {
	id := id.MakeId(rtype.Name())
	if was, exists := reg[id]; exists && was != rtype {
		err = errutil.New("has conflicting names, id:", id, "was:", was, "type:", rtype)
	} else {
		reg[id] = rtype
	}
	return
}
