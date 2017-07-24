package group

import (
	"github.com/ionous/iffy/rt"
)

type GroupList struct {
	Order   int // a sorting key
	Objects []rt.Object
}

func (l *GroupList) Empty() bool {
	return len(l.Objects) == 0
}

func (l *GroupList) Append(obj rt.Object) {
	l.Objects = append(l.Objects, obj)
}

type PendingGroups map[Key]GroupList

func (groups PendingGroups) Add(k Key, obj rt.Object) {
	if list := groups[k]; !list.Empty() {
		list.Append(obj)
		groups[k] = list
	} else {
		groups[k] = GroupList{
			len(groups),
			[]rt.Object{obj},
		}
	}
}

func (groups PendingGroups) Sort() []Collection {
	// we need to walk groups in some consistent order.
	sorted := make([]Collection, len(groups))
	for k, list := range groups {
		if i, objs := list.Order, list.Objects; len(objs) > 1 {
			sorted[i] = Collection{k, objs}
		}
	}
	return sorted
}
