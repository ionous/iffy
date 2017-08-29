package ops

import (
	r "reflect"
)

// Target handles the differences between structs and constructors.
// Its a subset of reflect.Value's functions.
type Target interface {
	Type() r.Type
	Field(int) r.Value
	FieldByName(string) r.Value
	Addr() r.Value
}
