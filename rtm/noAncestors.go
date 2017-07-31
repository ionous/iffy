package rtm

import (
	"github.com/ionous/iffy/rt"
)

// NoAncestors provides a default for rtm if one is not set.
type NoAncestors struct{}

func (NoAncestors) GetAncestors(rt.Runtime, rt.Object) (rt.ObjectStream, error) {
	return rt.EmptyObjects{}, nil
}
