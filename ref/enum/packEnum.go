package enum

import (
	"github.com/ionous/iffy/ref/coerce"
	r "reflect"
)

// Pack changes a named choice into an indexed choice.
// Requires that dst be the value of an enumerated type.
func Pack(dst, src r.Value) (okay bool) {
	if choices := Enumerate(dst.Type()); len(choices) > 0 {
		choice := src.String()
		if i, ok := StringToIndex(choice, choices); ok {
			okay = coerce.Value(dst, r.ValueOf(i)) == nil
		}
	}
	return
}

// Unpack changes an indexed choice into a named choice.
// Requires that src be the value of an enumerated type.
func Unpack(dst, src r.Value) (okay bool) {
	if choices := Enumerate(src.Type()); len(choices) > 0 {
		if c, ok := IndexToString(int(src.Int()), choices); ok {
			okay = coerce.Value(dst, r.ValueOf(c)) == nil
		}
	}
	return
}
