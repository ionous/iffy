package group

import (
	"github.com/ionous/iffy/rt"
)

type GroupList struct {
	Order   int // a sorting key
	Objects []rt.Object
}

func (l GroupList) Len() int {
	return len(l.Objects)
}

type ObjectGroups map[Key]GroupList

type PendingGroups struct {
	Groups  ObjectGroups
	Objects GroupedObjects
}

func (p *PendingGroups) Add(key Key, obj rt.Object) {
	// add all objects to the list
	p.Objects = append(p.Objects, GroupedObject{key, obj})
	// but only valid groups to the map
	if key.IsValid() {
		if list, ok := p.Groups[key]; ok {
			list.Objects = append(list.Objects, obj)
			p.Groups[key] = list
		} else {
			order := len(p.Groups)
			p.Groups[key] = GroupList{
				order,
				[]rt.Object{obj},
			}
		}
	}
}

func (p *PendingGroups) Sort() (Collections, []rt.Object) {
	// we need to walk p in some consistent order.
	sorted := make(Collections, len(p.Groups))
	for key, list := range p.Groups {
		if i, objs := list.Order, list.Objects; list.Len() > 1 {
			sorted[i] = Collection{key, objs}
		}
	}
	return sorted, p.Objects.Distill(p.Groups)
}
