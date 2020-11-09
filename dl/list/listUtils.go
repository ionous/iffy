package list

import (
	"github.com/ionous/iffy/affine"
	g "github.com/ionous/iffy/rt/generic"
)

func isList(v g.Value) bool {
	elAffinity := affine.Element(v.Affinity())
	return len(elAffinity) > 0
}

// increase the size of vs by amt, return the previous len of vs
func grow(vs g.Value, amt int) (retSlice g.Value, retSize int, err error) {
	if oldCnt, e := vs.GetLen(); e != nil {
		err = e
	} else if v, e := vs.Resize(oldCnt + amt); e != nil {
		err = e
	} else {
		retSlice, retSize = v, oldCnt
	}
	return
}
