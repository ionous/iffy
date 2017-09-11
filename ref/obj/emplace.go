package obj

import (
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
)

// Emplace wraps the passed value as an anonymous object.
func Emplace(i interface{}) rt.Object {
	rval, e := unique.ValuePtr(i)
	if e != nil {
		panic(e)
	}
	return RefObject{value: rval}
}
