package generic

import "github.com/ionous/iffy/affine"

func AreEqualTypes(a, b Value) bool {
	return areEqualTypes(a.Affinity(), a.Type(), b.Affinity(), b.Type())
}

func areEqualTypes(fa affine.Affinity, ft string, va affine.Affinity, vt string) (okay bool) {
	if fa == va {
		recordLike := fa == affine.Object || fa == affine.Record || fa == affine.RecordList
		okay = !recordLike || vt == ft
	}
	return
}
