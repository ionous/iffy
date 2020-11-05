package generic

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

type StringSlice struct {
	Nothing
	vals []string
}

func NewStringSlice(vals []string) *StringSlice {
	return &StringSlice{vals: vals}
}

func (l *StringSlice) Affinity() affine.Affinity { return affine.TextList }
func (l *StringSlice) Type() string              { return "string" }
func (l *StringSlice) GetTextList() (ret []string, _ error) {
	ret = l.vals
	return
}

// GetLen returns the number of elements in the underlying value if it's a slice,
// otherwise this returns an error.
func (l *StringSlice) GetLen() (ret int, _ error) {
	ret = len(l.vals)
	return

}

// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
// otherwise this returns an error.
func (l *StringSlice) GetIndex(i int) (ret rt.Value, err error) {
	vs := l.vals
	if cnt := len(vs); i < 0 || i >= cnt {
		err = rt.OutOfRange{i, cnt}
	} else {
		ret = NewString(vs[i])
	}
	return
}
