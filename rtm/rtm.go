package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

type Rtm struct {
	ref.Classes
	ref.Objects
	ref.Relations
	ScopeStack
	OutputStack
	Randomizer
}

func NewRtm(classes ref.Classes, objects ref.Objects, rel ref.Relations) *Rtm {
	rtm := &Rtm{Classes: classes, Objects: objects, Relations: rel}
	rtm.PushScope(scope.ModelFinder(rtm))
	rtm.Randomizer.Reset(1) // FIX: time?
	return rtm
}

func (rtm *Rtm) NewObject(class string) (ret rt.Object, err error) {
	id := id.MakeId(class)
	if cls, ok := rtm.Classes[id]; !ok {
		err = errutil.New("no such class", class)
	} else {
		ret = cls.NewObject()
	}
	return
}
