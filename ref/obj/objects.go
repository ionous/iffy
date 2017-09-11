package obj

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt"
)

//
type ObjectMap map[ident.Id]RefObject

// GetObject is compatible with rt.Runtime. The map can also be used directly.
func (or ObjectMap) GetObject(name string) (ret rt.Object, okay bool) {
	id := ident.IdOf(name)
	ret, okay = or[id]
	return
}
