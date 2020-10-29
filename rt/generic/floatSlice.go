package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

type FloatSlice struct {
	Nothing
	Values []float64
}

func (l *FloatSlice) Affinity() affine.Affinity { return affine.NumList }
func (l *FloatSlice) Type() string              { return "[]float64" }
func (l *FloatSlice) GetNumList() (ret []float64, _ error) {
	ret = l.Values
	return
}

// GetLen returns the number of elements in the underlying value if it's a slice,
// otherwise this returns an error.
func (l *FloatSlice) GetLen() (ret int, _ error) {
	ret = len(l.Values)
	return

}

// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
// otherwise this returns an error.
func (l *FloatSlice) GetIndex(i int) (ret rt.Value, err error) {
	if vs := l.Values; i < 0 || i >= len(vs) {
		err = OutOfRange(i)
	} else {
		ret = &Float{Value: vs[i]}
	}
	return
}

type OutOfRange int

func (e OutOfRange) Error() string {
	return errutil.Sprint("out of range", int(e))
}
