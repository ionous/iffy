package play

import (
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/rt"
)

type ParentChildAncestry struct{}

func (ParentChildAncestry) GetAncestors(run rt.Runtime, child rt.Object) (rt.ObjectStream, error) {
	return locate.GetAncestors(run, child)
}
