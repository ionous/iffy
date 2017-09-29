package stream

import (
	r "reflect"
)

// ValueError holds either a valid value, or an error regarding why the value couldn't be retrieved.
// Useful for iterating over a series of iffy calls via a linq-style function iterator
// ( because since iffy has so many explicit error returns. )
type ValueError struct {
	Value interface{}
	Error error
}

// Value returns a ValueError pair with the passed value, and no error.
// Useful when creating ValueError iterator functions compatible with linq;
// the second returned value therefore is always true.
func Value(v interface{}) (ValueError, bool) {
	return ValueError{Value: v}, true
}

// Error returns a ValueError pair with the passed error, and no value.
// Useful when creating ValueError iterator functions compatible with linq;
// the second returned value therefore is always true.
func Error(e error) (ValueError, bool) {
	return ValueError{Error: e}, true
}

// FromList turns the source list into a ValueError iterator compatible with linq.
// linq functions normally return Query, but iffy's streams are not restartable,
// so an iterator is a better match.
func FromList(source interface{}) func() (interface{}, bool) {
	src, i := r.ValueOf(source), 0
	return func() (ret interface{}, okay bool) {
		if i < src.Len() {
			el := src.Index(i)
			ret, okay = Value(el.Interface())
			i += 1
		}
		return
	}
}
