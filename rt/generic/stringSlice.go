package generic

import "github.com/ionous/iffy/rt"

type StringSlice struct {
	Nothing
	Values []string
}

// GetLen returns the number of elements in the underlying value if it's a slice,
// otherwise this returns an error.
func (l *StringSlice) GetLen(rt.Runtime) (ret int, _ error) {
	ret = len(l.Values)
	return

}

// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
// otherwise this returns an error.
func (l *StringSlice) GetIndex(_ rt.Runtime, i int) (ret rt.Value, err error) {
	if vs := l.Values; i < 0 || i >= len(vs) {
		err = OutOfRange(i)
	} else {
		ret = &String{Value: vs[i]}
	}
	return
}
