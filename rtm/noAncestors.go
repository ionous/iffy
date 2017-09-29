package rtm

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// NoAncestors provides a default for rtm if one is not set.
type NoAncestors struct{}

func (NoAncestors) GetAncestors(rt.Runtime, rt.Object) (rt.ObjectStream, error) {
	return stream.NewObjectStream(func() (ret interface{}, okay bool) {
		return
	}), nil
}
