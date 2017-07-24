package group

import (
	"github.com/ionous/iffy/rt"
)

// GroupedObject reverses the mapping between object and group.
type GroupedObject struct {
	Group  *Key
	Object rt.Object
}

// GroupedObjects holds the reverse mapping of PendingGroups.
// We accumulate every object that could be grouped or ungrouped and then distill the list into the ungrouped objects.
type GroupedObjects []GroupedObject

// Distill returns a list of ungrouped objects
func (r GroupedObjects) Distill(groups PendingGroups) (ret []rt.Object) {
	for _, x := range r {
		if x.Group == nil || len(groups[*x.Group].Objects) == 1 {
			ret = append(ret, x.Object)
		}
	}
	return
}
