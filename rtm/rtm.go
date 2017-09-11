package rtm

import (
	"github.com/ionous/iffy/event"
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
	obj.ObjectMap
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
	return rtm.GetObject(name)
}

// GetValue sets the value of the passed pointer to the value of the named property in the passed object.
func (rtm *Rtm) GetValue(src rt.Object, name string, pv interface{}) error {
	return obj.UnpackValue(rtm.ObjectMap, src, name, pv)
}

// SetValue sets the named property in the passed object to the value.
func (rtm *Rtm) SetValue(dst rt.Object, name string, v interface{}) error {
	return obj.PackValue(rtm.ObjectMap, dst, name, v)
}
