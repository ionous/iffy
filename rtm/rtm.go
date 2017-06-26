package rtm

import (
	"github.com/ionous/iffy/ref"
)

type Rtm struct {
	ref.Model
	ScopeStack
	OutputStack
}

func NewRtm(model ref.Model) *Rtm {
	rtm := &Rtm{Model: model}
	rtm.PushScope(ModelFinder{model})
	return rtm
}
