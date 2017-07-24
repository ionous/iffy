package group

import (
	"github.com/ionous/iffy/rt"
)

// Key partitions objects into groups.
type Key struct {
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
}

// Collection contains a group of objects.
type Collection struct {
	Key
	Objects []rt.Object
}

// ObjectGrouping defines how objects in groups should display.
type ObjectGrouping int

//go:generate stringer -type=ObjectGrouping
const (
	WithoutArticles ObjectGrouping = iota
	WithArticles
	WithoutObjects
)

// GroupTogether executes a pattern to collect objects.
// FIX: ideally the member here would be "Key", but we need op/spec to handle aggregates.
type GroupTogether struct {
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
	Target         rt.Object
}

// PrintGroup executes a pattern to print a collection of objects.
// FIX: ideally the member here would be "Key", but we need op/spec to handle aggregates.
type PrintGroup struct {
	Label          string
	Innumerable    bool
	ObjectGrouping ObjectGrouping
	Objects        []rt.Object
}
