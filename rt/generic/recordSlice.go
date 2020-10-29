package generic

// import (
// 	"github.com/ionous/iffy/affine"
// 	"github.com/ionous/iffy/rt"
// )

// type RecordSlice struct {
// 	Nothing
// 	kind   string
// 	values []fieldData
// }

// var _ rt.Value = (*RecordSlice)(nil) // ensure compatibility

// func NewRecordSlice(kind string) *RecordSlice {
// 	return &Record{kind: kind, fields: fields}
// }

// func (l *RecordSlice) Affinity() affine.Affinity { return affine.RecordList }
// func (l *RecordSlice) Type() string              { return "[]" + l.kind }

// // GetLen returns the number of elements in the underlying value if it's a slice,
// // otherwise this returns an error.
// func (l *RecordSlice) GetLen() (ret int, _ error) {
// 	ret = len(l.Values)
// 	return

// }

// // GetIndex returns the nth element of the underlying slice, where 0 is the first value;
// // otherwise this returns an error.
// func (l *RecordSlice) GetIndex(i int) (ret rt.Value, err error) {
// 	if vs := l.Values; i < 0 || i >= len(vs) {
// 		err = OutOfRange(i)
// 	} else {
// 		ret = &Record{kind: l.kind, fields: vs[i]}
// 	}
// 	return
// }
