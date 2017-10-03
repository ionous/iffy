package group

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec"
)

// Key partitions objects into groups.
// A label can replace the listing of individual objects in the group, or it can collect them into a bracket ( depending on the ObjectGrouping setting ). Except if innumerable, the number of objects in the group will be printed just before the label.
type Key struct {
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
}

// IsValid returns true for keys with either a label or a display of individual objects
func (k Key) IsValid() bool {
	return len(k.Label) > 0 || k.ObjectGrouping != WithoutObjects
}

// Collection contains a group of objects.
type Collection struct {
	Key
	Objects []rt.Object
}

type Collections []Collection

// ObjectGrouping defines how objects in groups should display.
type ObjectGrouping int

//go:generate stringer -type=ObjectGrouping
const (
	WithoutObjects ObjectGrouping = iota
	WithoutArticles
	WithArticles
)

// GroupTogether executes a pattern to collect objects.
type GroupTogether struct {
	// FIX: ideally the member here would be "Key", but we need op/spec to handle aggregates.
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
	Target         rt.Object
}

// PrintGroup executes a pattern to print a collection of objects.
type PrintGroup struct {
	// FIX: ideally the member here would be "Key", but we need op/spec to handle aggregates.
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
	Objects        []rt.Object
}

func GroupRules(c spec.Block) {
	// all unnamed objects go into a group w/ text the plural of the class name.
	if c.Cmd("run rule", "group together").Begin() {
		c.Param("if").Cmd("includes", c.Cmd("get", c.Cmd("get", "@", "target"), "name"), "#")
		if c.Param("decide").Cmds().Begin() {
			c.Cmd("set text", "@", "label", c.Cmd("pluralize", c.Cmd("class name", c.Cmd("get", "@", "target"))))
			c.End()
		}
		c.End()
	}
	if c.Cmd("run rule", "print group").Begin() {
		if c.Param("decide").Cmds().Begin() {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					if c.Cmd("choose", c.Cmd("is empty", c.Cmd("get", "@", "label"))).Begin() {
						// no label: then the group is an inline list of things ( and no brackets. )
						if c.Param("true").Cmds().Begin() {
							if c.Cmd("print objects", c.Cmd("get", "@", "objects")).Begin() {
								c.Param("articles").Cmd("get", "@", "with articles")
								c.End() // print objects
							}
							c.End() // true, no label
						}
						// the label is not empty: then the group is a block of things.
						if c.Param("false").Cmds().Begin() {
							// before the label, possibly write the number of objects:
							if c.Cmd("choose", c.Cmd("get", "@", "innumerable")).Begin() {
								if c.Param("false").Cmds().Begin() {
									c.Cmd("print num word", c.Cmd("len", c.Cmd("get", "@", "objects")))
									c.End()
								}
								c.End()
							}
							// now the label:
							c.Cmd("say", c.Cmd("get", "@", "label"))
							// after the label, possibly write the objects:
							if c.Cmd("choose", c.Cmd("get", "@", "without objects")).Begin() {
								if c.Param("false").Cmds().Begin() {
									if c.Cmd("choose", c.Cmd("get", "@", "innumerable")).Begin() {
										// if they are not innumerable, they are numerable.
										// if they are numerable, then they got a number... in front of some brackets.
										if c.Param("false").Cmds().Begin() {
											if c.Cmd("print bracket").Begin() {
												if c.Cmds().Begin() {
													c.Cmd("print objects", c.Cmd("get", "@", "objects"),
														c.Param("articles").Cmd("get", "@", "with articles"))
													c.End()
												}
												c.End()
											}
											c.End()
										}
										if c.Param("true").Cmds().Begin() {
											c.Cmd("print objects", c.Cmd("get", "@", "objects"),
												c.Param("articles").Cmd("get", "@", "with articles"))
											c.End()
										}
										c.End()
									}
									c.End() // false, not without objects
								}
								c.End() // choose objects
							}
							c.End() // false, has label
						}
						c.End() // end of choose
					}
					c.End() // end of line statements
				}
				c.End() // end of print
			}
			c.End() // end of decide
		}
		c.End() // end of pattern
	}
}
