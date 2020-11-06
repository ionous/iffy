package generic

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/chain"
)

func SliceFloats(vs []float64) rt.Iterator {
	return chain.SliceIt(len(vs), func(i int) (rt.Value, error) {
		return NewFloat(vs[i]), nil
	})
}

func SliceStrings(vs []string) rt.Iterator {
	return chain.SliceIt(len(vs), func(i int) (rt.Value, error) {
		return NewString(vs[i]), nil
	})
}
