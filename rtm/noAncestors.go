package rtm

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

// NoAncestors provides a default for rtm if one is not set.
type NoAncestors struct{}

func (NoAncestors) GetAncestors(rt.Object) (rt.ObjectStream, error) {
	return rt.EmptyObjects{}, nil
}
