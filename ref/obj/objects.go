package obj

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

//
type ObjectMap map[ident.Id]RefObject

// Emplace wraps the passed value as an anonymous object.
// Compatible with rt.Runtime.
func (or ObjectMap) Emplace(i interface{}) rt.Object {
	return Emplace(i)
}

// GetObject is compatible with rt.Runtime. The map can also be used directly.
func (or ObjectMap) GetObject(name string) (ret rt.Object, okay bool) {
	id := ident.IdOf(name)
	ret, okay = or[id]
	return
}
