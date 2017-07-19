package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	r "reflect"
)

// PrintName executes a pattern to print the target's name.
// The standard rules print the "printed name" property of the target,
// or the object name ( if the target lacks a "printed name" ),
// or the object's class name ( for unnamed objects. )
// A "printed name" can change during the course of play; object names never change.
type PrintName struct {
	Target *Kind
}

// PrintPluralName executes a pattern to print the plural of the target's name.
// The standard rules print the target's "printed plural name",
// or, if the target lacks that property, the plural of the "print name" pattern.
// It uses the runtime's pluralization table, or if needed, automated pluralization.
type PrintPluralName struct {
	Target *Kind
}

// Group defines a collection of objects.
type Group struct {
	Text           string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
}

// GroupTogether executes a pattern to accumulate groupings.
// FIX: ideally the member here would be "Group", but we need op/spec to handle aggregates.
type GroupTogether struct {
	Text           string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
	Target         rt.Object // FIX: what is this supposed to be?
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
	Text           string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
	Objects        []rt.Object // FIX: what is this supposed to be?
}

// PrintNondescriptObjects commands the runtime to print a bunch of objects, in groups if possible.
type PrintNondescriptObjects struct {
	Objects rt.ObjListEval
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

func (p *PrintNondescriptObjects) Execute(run rt.Runtime) (err error) {
	if groups, e := makeGroups(run, p.Objects); e != nil {
		err = e
	} else {
		outer := printer.AndSeparator()
		if e := rt.Write(run, outer, func() (err error) {
			// sort the map
			type Sorted struct {
				Group
				objs []rt.Object
			}
			sorted := make([]Sorted, len(groups.Grouped))
			for group, list := range groups.Grouped {
				s := &(sorted[list.Order])
				s.Group = group
				s.objs = list.Objects
			}
			for _, group := range sorted {
				if printGroup, e := run.Emplace(&PrintGroup{group.Text, group.Innumerable, group.ObjectGrouping, group.objs}); e != nil {
					err = e
					break
				} else {
					inner := printer.AndSeparator()
					if e := rt.Write(run, inner, func() error {
						_, e := run.ExecuteMatching(printGroup)
						_, e = pat.Found(e)
						return e
					}); e != nil {
						err = e
						break
					} else if e := printList(run, group.objs); e != nil {
						err = e
						break
					}
					// flush the group memebers to the output
					if _, e := inner.WriteTo(outer); e != nil {
						err = e
						break
					}
				}
			}
			if err == nil {
				err = printList(run, groups.Ungrouped)
			}
			return
		}); e != nil {
			err = e
		} else {
			// flush the group phrases to the output
			if _, e := outer.WriteTo(run); e != nil {
				err = e
			}
		}
	}
	return
}

func printList(run rt.Runtime, objs []rt.Object) (err error) {
	for _, obj := range objs {
		if e := printName(run, obj); e != nil {
			err = e
			break
		}
	}
	return
}

// makeGroups breaks the passed stream of objects into separated units.
func makeGroups(run rt.Runtime, ol rt.ObjListEval) (groups Groups, err error) {
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
				} else {
					ran, e := run.ExecuteMatching(grouped)
					if _, e := pat.Found(e); e != nil {
						err = e
						break
					} else if !ran {
						groups.Ungrouped = append(groups.Ungrouped, obj)
					} else {
						// add the object to the group
						group := Group{group.Text, group.Innumerable, group.ObjectGrouping}
						if list := groups.Grouped[group]; !list.Empty() {
							list.Append(obj)
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
	}
	return
}

// FIX: this is patently ridiculous.
// issue: i cant set an object reference from an object
// why? in part b/c theres no "base class"
// it would be **alot** simpler if the * was an ident.Id
// we'd still have "emplace" -- you could maybe someday make it static -- thatd be tons better.
func printName(run rt.Runtime, obj rt.Object) (err error) {
	var kind *Kind
	if src, ok := obj.(*ref.RefObject); !ok {
		err = errutil.Fmt("unknown object %T", obj)
	} else if e := ref.Upcast(src.Value().Addr(), func(ptr r.Value) (okay bool) {
		kind, okay = ptr.Interface().(*Kind)
		return
	}); e != nil {
		err = e
	} else if printName, e := run.Emplace(&PrintName{kind}); e != nil {
		err = e
	} else {
		_, err = run.ExecuteMatching(printName)
	}
	return
}
