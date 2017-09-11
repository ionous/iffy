package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/rel"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"io"
	r "reflect"
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

// GetValue sets the value of the passed pointer to the value of the named property in the passed object.
func (rtm *Rtm) GetValue(obj rt.Object, name string, pv interface{}) (err error) {
	if pdst := r.ValueOf(pv); pdst.Kind() != r.Ptr {
		err = errutil.New(obj, name, "expected pointer outvalue", pdst.Type())
	} else if p, ok := obj.Property(name); !ok {
		err = errutil.New(obj, name, "unknown property")
	} else {
		dst := pdst.Elem()
		src := r.ValueOf(p.Value())
		if e := rtm.pack(dst, src); e != nil {
			err = errutil.New(obj, name, "cant unpack", dst.Type(), "from", src.Type(), "because", e)
		}
	}
	return
}

// SetValue sets the named property in the passed object to the value.
func (rtm *Rtm) SetValue(obj rt.Object, name string, v interface{}) (err error) {
	if p, ok := obj.Property(name); !ok {
		err = errutil.New(obj, name, "unknown property")
	} else {
		dst := r.New(p.Type()).Elem() // create a new destination for the value.
		src := r.ValueOf(v)
		if e := rtm.pack(dst, src); e != nil {
			err = errutil.New(obj, name, "cant pack", dst.Type(), "from", src.Type(), "because", e)
		} else {
			err = p.SetValue(dst.Interface())
		}
	}
	return
}
