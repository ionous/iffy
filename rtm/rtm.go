package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"io"
)

type Rtm struct {
	unique.Types
	ref.ObjectMap
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

// FindObject implements rt.ObjectFinder.
func (rtm *Rtm) FindObject(name string) (rt.Object, bool) {
	return rtm.GetObject(name)
}

// GetValue sets the value of the passed pointer to the value of the named property in the passed object.
func (rtm *Rtm) GetValue(obj rt.Object, name string, pv interface{}) (err error) {
	if p, ok := obj.Property(name); !ok {
		err = errutil.Fmt("unknown property %s.%s", obj, name)
	} else {
		err = p.GetValue(pv)
	}
	return
}

// SetValue sets the named property in the passed object to the value.
func (rtm *Rtm) SetValue(obj rt.Object, name string, v interface{}) (err error) {
	if p, ok := obj.Property(name); !ok {
		err = errutil.Fmt("unknown property %s.%s", obj, name)
	} else {
		err = p.SetValue(v)
	}
	return
}
