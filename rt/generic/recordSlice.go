package generic

import (
	"github.com/ionous/iffy/affine"
)

type RecordSlice struct {
	Nothing
	kind   *Kind
	values []*Record
}

var _ Value = (*RecordSlice)(nil) // ensure compatibility

func (l *RecordSlice) Affinity() affine.Affinity { return affine.RecordList }
func (l *RecordSlice) Type() string              { return l.kind.name }
func (l *RecordSlice) GetRecordList() (ret []*Record, _ error) {
	ret = l.values
	return
}

// GetLen returns the number of elements in the underlying value if it's a slice,
// otherwise this returns an error.
func (l *RecordSlice) GetLen() (ret int, _ error) {
	ret = len(l.values)
	return
}

// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
// otherwise this returns an error.
func (l *RecordSlice) GetIndex(i int) (ret Value, err error) {
	if vs := l.values; i < 0 {
		err = Underflow{i, 0}
	} else if cnt := len(vs); i >= cnt {
		err = Overflow{i, cnt}
	} else {
		ret = vs[i]
	}
	return
}
