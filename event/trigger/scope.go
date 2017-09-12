package trigger

import (
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
)

type Context struct {
	run    rt.Runtime
	events event.EventMap
}

type ObjectNoun struct{ rt.Object }

type Probe struct {
	Run    rt.Runtime
	Search func(rt.Runtime, parser.NounVisitor) bool
}

// ex. "held"
func (m Context) GetPlayerScope(n string) (ret parser.Scope, err error) {
	if len(n) > 0 {
		panic("not implemented")
	}
	ret = Probe{m.run, SearchPlayerScope}
	return
}

func (m Context) GetObjectScope(id ident.Id) (ret parser.Scope, err error) {
	panic("not implemented")
}

func (m Context) IsPlural(word string) bool {
	return word != lang.Singularize(word)
}

func SearchPlayerScope(run rt.Runtime, nv parser.NounVisitor) (ret bool) {
	// for a bunch of objects in the runtime.
	rtm := run.(*rtm.Rtm)
	for _, v := range rtm.Objects {
		if nv(ObjectNoun{v}) {
			ret = true
			break
		}
	}
	return
}

func (on ObjectNoun) HasName(string) bool {
	// not completely sure where to do this yet
	panic("split object name into parts")
}
func (on ObjectNoun) HasPlural(string) bool {
	// ex. "take apples"
	// we can use pluralized class names  --
	// its possible we might want to check printed plural name as well.
	panic("test class plural")
}
func (on ObjectNoun) HasClass(string) bool {
	// ex. only apply this action to "things"
	panic("test for exact class name")
}
func (on ObjectNoun) HasAttribute(string) bool {
	panic("test for object state")
}

func (w Probe) SearchScope(nv parser.NounVisitor) bool {
	return w.Search(w.Run, nv)
}
