package generic

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

type StringSlice struct {
	Nothing
	Values []string
}

func (l *StringSlice) Affinity() affine.Affinity { return affine.TextList }
func (l *StringSlice) Type() string              { return "[]string" }
func (l *StringSlice) GetTextList() (ret []string, _ error) {
	ret = l.Values
	return
}

// GetLen returns the number of elements in the underlying value if it's a slice,
// otherwise this returns an error.
func (l *StringSlice) GetLen() (ret int, _ error) {
	ret = len(l.Values)
	return

}

// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
// otherwise this returns an error.
func (l *StringSlice) GetIndex(i int) (ret rt.Value, err error) {
	if vs := l.Values; i < 0 || i >= len(vs) {
		err = OutOfRange(i)
	} else {
		ret = &String{Value: vs[i]}
	}
	return
}
