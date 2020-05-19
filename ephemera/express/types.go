package express

import (
	r "reflect"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

var typeNumEval = r.TypeOf((*rt.NumberEval)(nil)).Elem()
var typeTextEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var compareNum = r.TypeOf((*core.CompareNum)(nil)).Elem()
var compareText = r.TypeOf((*core.CompareText)(nil)).Elem()

func implements(a, b r.Value, t r.Type) bool {
	return a.Type().Implements(t) && b.Type().Implements(t)
}
