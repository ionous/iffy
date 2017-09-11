package rtm

import (
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/rel"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"io"
)

type Rtm struct {
	unique.Types
	Objects obj.ObjectMap
	rel.Relations
	io.Writer
	Events event.EventMap
	Randomizer
	rt.Ancestors
	pat.Patterns
	Plurals
	parser.Scanner
}

// GetPatterns mainly for testing.
func (rtm *Rtm) GetPatterns() *pat.Patterns {
	return &rtm.Patterns
}

// FindObject implements rt.ObjectFinder.
func (rtm *Rtm) FindObject(name string) (rt.Object, bool) {
	return rtm.Objects.GetObject(name)
}

// GetValue from the map
func (rtm *Rtm) GetObject(name string) (rt.Object, bool) {
	return rtm.Objects.GetObject(name)
}

func (rtm *Rtm) Emplace(i interface{}) rt.Object {
	return obj.MakeObject(ident.None(), i, rtm)
}
