package rtm

import (
	"github.com/ionous/iffy/ref"
)

type Rtm struct {
	ref.Model
	ScopeStack
}

func NewRtm(model ref.Model) *Rtm {
	return &Rtm{Model: model}
}
