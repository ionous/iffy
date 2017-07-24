package group

import (
	"github.com/ionous/iffy/rt"
)

// GroupedObject reverses the mapping between object and group.
type GroupedObject struct {
	Key    Key
	Object rt.Object
}

// GroupedObjects holds the reverse mapping of PendingGroups.
// We accumulate every object that could be grouped or ungrouped and then distill the list into the ungrouped objects.
type GroupedObjects []GroupedObject

// Distill returns a list of ungrouped objects
func (r GroupedObjects) Distill(groups ObjectGroups) (ret []rt.Object) {
	for _, x := range r {
		if groups[x.Key].Len() < 2 {
			ret = append(ret, x.Object)
		}
	}
	return
}
