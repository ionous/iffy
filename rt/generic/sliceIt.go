package generic

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/chain"
)

func SliceFloats(vs []float64) rt.Iterator {
	return chain.SliceIt(len(vs), func(i int) (ret rt.Value, err error) {
		ret = FloatOf(vs[i])
		return
	})
}

func SliceStrings(vs []string) rt.Iterator {
	return chain.SliceIt(len(vs), func(i int) (ret rt.Value, err error) {
		ret = StringOf(vs[i])
		return
	})
}
