package ref

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
)

// Relations maps ids to RefReleation.
// Compatible with unique.TypeRegistry.
type Relations map[string]*RefRelation

func (reg Relations) GetRelation(name string) (ret rt.Relation, okay bool) {
	id := id.MakeId(name)
	if ref, ok := reg[id]; ok {
		ret, okay = ref, true
	}
	return
}
