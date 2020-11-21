package list

import (
	"github.com/ionous/iffy/affine"
	g "github.com/ionous/iffy/rt/generic"
)

func isList(v g.Value) bool {
	return affine.IsList(v.Affinity())
}
