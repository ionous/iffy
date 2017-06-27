package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
	r "reflect"
)

type Rtm struct {
	ref.Classes
	ref.Objects
	ScopeStack
	OutputStack
	Randomizer
}

func NewRtm(classes ref.Classes, objects ref.Objects) *Rtm {
	rtm := &Rtm{Classes: classes, Objects: objects}
	rtm.PushScope(scope.ModelFinder(rtm))
	rtm.Randomizer.Reset(1) // FIX: time?
	return rtm
}

func (rtm *Rtm) NewObject(class string) (ret rt.Object, err error) {
	id := id.MakeId(class)
	if cls, ok := rtm.Classes[id]; !ok {
		err = errutil.New("no such class", class)
	} else {
		inst := r.New(cls.Type())
		ret = ref.NewInst(cls, inst.Elem())
	}
	return
}
