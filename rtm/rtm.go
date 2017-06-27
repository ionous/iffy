package rtm

import (
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt/scope"
)

type Rtm struct {
	ref.Model
	ScopeStack
	OutputStack
	Randomizer
}

func NewRtm(model ref.Model) *Rtm {
	rtm := &Rtm{Model: model}
	rtm.PushScope(scope.ModelFinder(model))
	rtm.Randomizer.Reset(1) // FIX: time?
	return rtm
}
