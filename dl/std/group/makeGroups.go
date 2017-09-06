package group

import (
	"github.com/ionous/iffy/rt"
)

// MakeGroups breaks the passed stream of objects into separated units.
// Note: not all returned groups will have objects.
func MakeGroups(run rt.Runtime, ol rt.ObjListEval) (groups Collections, ungrouped []rt.Object, err error) {
	if os, e := ol.GetObjectStream(run); e != nil {
		err = e
	} else {
		pending := PendingGroups{make(ObjectGroups), nil}
		//
		for os.HasNext() {
			if obj, e := os.GetNext(); e != nil {
				err = e
				break
			} else {
				// find the desired group for this object.
				group := GroupTogether{Target: obj}
				grouped :=run.Emplace(&group)
				if e := run.ExecuteMatching(run, grouped); e != nil {
					err = e
					break
				} else {
					// if nothing was set, then the key is invalid, and the object is considered ungrouped.
					key := Key{group.Label, group.Innumerable, group.ObjectGrouping}
					pending.Add(key, obj)
				}
			
		}
		if err == nil {
			groups, ungrouped = pending.Sort()
		}
	}
	return
}
