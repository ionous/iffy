package group

import (
	"github.com/ionous/iffy/rt"
)

// MakeGroups breaks the passed stream of objects into separated units.
// Note: not all returned groups will have objects.
func MakeGroups(run rt.Runtime, ol rt.ObjListEval) (groups []Collection, ungrouped []rt.Object, err error) {
	if os, e := ol.GetObjectStream(run); e != nil {
		err = e
	} else {
		var ungroup GroupedObjects
		pending := make(PendingGroups)
		//
		for os.HasNext() {
			if obj, e := os.GetNext(); e != nil {
				err = e
				break
			} else {
				// find the desired group for this object.
				group := GroupTogether{Target: obj}
				if grouped, e := run.Emplace(&group); e != nil {
					err = e
					break
				} else if ran, e := run.ExecuteMatching(grouped); e != nil {
					err = e
					break
				} else {
					var key *Key
					if ran {
						k := Key{group.Label, group.Innumerable, group.ObjectGrouping}
						pending.Add(k, obj)
						key = &k
					}
					ungroup = append(ungroup, GroupedObject{key, obj})
				}
			}
		}
		if err == nil {
			groups = pending.Sort()
			ungrouped = ungroup.Distill(pending)
		}
	}
	return
}
