package list

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

func isList(v rt.Value) bool {
	elAffinity := affine.Element(v.Affinity())
	return len(elAffinity) > 0
}

// increase the size of vs by amt, return the previous len of vs
func grow(vs rt.Value, amt int) (retSlice rt.Value, retSize int, err error) {
	if oldCnt, e := vs.GetLen(); e != nil {
		err = e
	} else if v, e := vs.Resize(oldCnt + amt); e != nil {
		err = e
	} else {
		retSlice, retSize = v, oldCnt
	}
	return
}
