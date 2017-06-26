package rtm

import (
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt/scope"
)

type Rtm struct {
	ref.Model
	ScopeStack
	OutputStack
}

func NewRtm(model ref.Model) *Rtm {
	rtm := &Rtm{Model: model}
	rtm.PushScope(scope.ModelFinder(model))
	return rtm
}
