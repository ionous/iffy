package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
)

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

type Groups map[Group]ObjectList
type ObjectList []rt.Object

func (p *PrintNondescriptObjects) Execute(run rt.Runtime) (err error) {
	if os, e := p.Objects.GetObjectStream(run); e != nil {
		err = e
	} else if groups, ungrouped, e := makeGroups(run, os); e != nil {
		err = e
	} else {
		outer := printer.AndSeparator()
		if e := rt.Write(run, outer, func() (err error) {
			for group, objs := range groups {
				if printGroup, e := run.Emplace(&PrintGroup{group.Text, group.Innumerable, group.ObjectGrouping, objs}); e != nil {
					err = e
					break
				} else {
					inner := printer.AndSeparator()
					if e := rt.Write(run, inner, func() error {
						_, e := run.ExecuteMatching(printGroup)
						return e
					}); e != nil {
						err = e
						break
					} else if e := printList(run, objs); e != nil {
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
				err = printList(run, ungrouped)
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
func makeGroups(run rt.Runtime, os rt.ObjectStream) (groups Groups, ungrouped ObjectList, err error) {
	groups = make(Groups)
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
				ungrouped = append(ungrouped, obj)
			} else {
				// add the object to the group
				group := Group{group.Text, group.Innumerable, group.ObjectGrouping}
				objs := groups[group]
				objs = append(objs, obj)
				groups[group] = objs
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
	if ref, ok := obj.(*ref.RefObject); !ok {
		err = errutil.Fmt("unknown object %T", obj)
	} else {
		ptr := ref.Value().Interface()
		if kind, ok := ptr.(*Kind); !ok {
			err = errutil.Fmt("unknown kind %T", ptr)
		} else if printName, e := run.Emplace(&PrintName{kind}); e != nil {
			err = e
		} else {
			_, err = run.ExecuteMatching(printName)
		}
	}
	return
}
