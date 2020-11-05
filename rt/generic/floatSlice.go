package generic

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

type FloatSlice struct {
	Nothing
	vals []float64
}

func NewFloatSlice(vals []float64) *FloatSlice {
	return &FloatSlice{vals: vals}
}

func (l *FloatSlice) Affinity() affine.Affinity { return affine.NumList }
func (l *FloatSlice) Type() string              { return "float64" }
func (l *FloatSlice) GetNumList() (ret []float64, _ error) {
	ret = l.vals
	return
}

// GetLen returns the number of elements in the underlying value if it's a slice,
// otherwise this returns an error.
func (l *FloatSlice) GetLen() (ret int, _ error) {
	ret = len(l.vals)
	return

}

// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
// otherwise this returns an error.
func (l *FloatSlice) GetIndex(i int) (ret rt.Value, err error) {
	if vs := l.vals; i < 0 {
		err = rt.OutOfRange{i, 0}
	} else if cnt := len(vs); i >= cnt {
		err = rt.OutOfRange{i, cnt}
	} else {
		ret = NewFloat(vs[i])
	}
	return
}
