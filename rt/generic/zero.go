package generic

import (
	"github.com/ionous/iffy/rt"
)

var (
	True      = BoolOf(true)
	False     = BoolOf(false)
	Zero      = FloatOf(0.0)
	ZeroList  = FloatsOf(nil)
	Empty     = StringOf("")
	EmptyList = StringsOf(nil)
)

const defaultType = "" // empty string

func must(v rt.Value, e error) rt.Value {
	if e != nil {
		panic(e)
	}
	return v
}
