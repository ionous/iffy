package rtm

import (
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"io"
)

type Rtm struct {
	unique.Types
	*ref.Objects
	ref.Relations
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

// GetClass implements rt.Model.
func (rtm *Rtm) GetClass(name string) (ret rt.Class, okay bool) {
	id := ident.IdOf(name)
	if cls, ok := rtm.Types[id]; ok {
		ret, okay = cls, ok
	}
	return
}

// FindObject implements rt.ObjectFinder.
func (rtm *Rtm) FindObject(name string) (rt.Object, bool) {
	return rtm.GetObject(name)
}

// GetValue sets the value of the passed pointer to the value of the named property in the passed object.
func (rtm *Rtm) GetValue(obj rt.Object, name string, pv interface{}) error {
	return obj.GetValue(name, pv)
}

// SetValue sets the named property in the passed object to the value.
func (rtm *Rtm) SetValue(obj rt.Object, name string, v interface{}) error {
	return obj.SetValue(name, v)
}
