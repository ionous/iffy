package rtm

import (
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt/scope"
)

type Rtm struct {
	*ref.Classes
	*ref.Objects
	*ref.Relations
	ScopeStack
	OutputStack
	Randomizer
}

func NewRtm(c *ref.Classes, o *ref.Objects, r *ref.Relations) *Rtm {
	rtm := &Rtm{
		Classes:   c,
		Objects:   o,
		Relations: r,
	}
	rtm.PushScope(scope.ModelFinder(rtm))
	rtm.Randomizer.Reset(1) // FIX: time?
	return rtm
}
