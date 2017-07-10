package pat

import (
	"github.com/ionous/errutil"
)

const NotFound errutil.Error = "pattern not found"

type Patterns struct {
	BoolMap
	NumberMap
	TextMap
	ObjectMap
	NumListMap
	TextListMap
	ObjListMap
}
