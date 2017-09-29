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
	writer io.Writer
	Events event.EventMap
	Randomizer
	rt.Ancestors
	Rules pat.Rulebook
	Plurals
	parser.Scanner
	ObjectScope
}

//
func (rtm *Rtm) GetObject(name string) (rt.Object, bool) {
	return rtm.Objects.GetObject(name)
}

//
func (rtm *Rtm) Emplace(i interface{}) rt.Object {
	return obj.MakeObject(ident.None(), i, rtm)
}

func (rtm *Rtm) Writer() io.Writer {
	return rtm.writer
}

func (rtm *Rtm) Write(b []byte) (int, error) {
	return rtm.writer.Write(b)
}

//
func (rtm *Rtm) SetWriter(w io.Writer) io.Writer {
	if w == nil {
		panic("push writer requires a valid object")
	}
	prev := rtm.writer
	rtm.writer = w
	return prev
}
