package core

import (
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// a utility, primarily used for testing, which allows values to be passed directly to commands which take parameters
type FromValue struct{ g.Value }

func (op *FromValue) GetEval() interface{} { return nil }

func (op *FromValue) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	ret = op.Value
	return
}
