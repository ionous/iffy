package std

import (
	"github.com/ionous/iffy/rt"
)

// Group defines a collection of objects.
type Group struct {
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
}

// GroupTogether executes a pattern to accumulate groupings.
// FIX: ideally the member here would be "Group", but we need op/spec to handle aggregates.
type GroupTogether struct {
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
	Target         rt.Object
}

type ObjectGrouping int

//go:generate stringer -type=ObjectGrouping
const (
	GroupWithoutArticles ObjectGrouping = iota
	GroupWithArticles
	GroupWithoutObjects
)

// PrintGroup executes a pattern to print a collection of objects.
// FIX: ideally the member here would be "Group", but we need op/spec to handle aggregates.
type PrintGroup struct {
	// FIX: op aggregates
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
	Objects        []rt.Object // FIX: what is this supposed to be?
}

type Groups struct {
	Grouped   map[Group]ObjectList
	Ungrouped []rt.Object
}

type ObjectList struct {
	Order   int
	Objects []rt.Object
}

func (l *ObjectList) Empty() bool {
	return len(l.Objects) == 0
}

func (l *ObjectList) Append(obj rt.Object) {
	l.Objects = append(l.Objects, obj)
}

// MakeGroups breaks the passed stream of objects into separated units.
func MakeGroups(run rt.Runtime, ol rt.ObjListEval) (groups Groups, err error) {
	if os, e := ol.GetObjectStream(run); e != nil {
		err = e
	} else {
		groups.Grouped = make(map[Group]ObjectList)
		//
		for os.HasNext() {
			if obj, e := os.GetNext(); e != nil {
				err = e
				break
			} else {
				group := GroupTogether{Target: obj}
				if grouped, e := run.Emplace(&group); e != nil {
					err = e
					break
				} else if ran, e := run.ExecuteMatching(grouped); e != nil {
					err = e
					break
				} else if !ran {
					groups.Ungrouped = append(groups.Ungrouped, obj)
				} else {
					// add the object to the group
					group := Group{group.Label, group.Innumerable, group.ObjectGrouping}
					if list := groups.Grouped[group]; !list.Empty() {
						list.Append(obj)
						groups.Grouped[group] = list
					} else {
						groups.Grouped[group] = ObjectList{
							len(groups.Grouped),
							[]rt.Object{obj},
						}
					}
				}
			}
		}
	}
	return
}
