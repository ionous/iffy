package rel

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt"
)

// Relations maps ids to RefReleation.
// Compatible with unique.TypeRegistry.
type Relations map[ident.Id]*RefRelation

func (reg Relations) GetRelation(name string) (ret rt.Relation, okay bool) {
	id := ident.IdOf(name)
	if ref, ok := reg[id]; ok {
		ret, okay = ref, true
	}
	return
}
