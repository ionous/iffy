package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type numEval struct{ eval rt.NumberEval }

// GetNumber, or error if the underlying value isn't a number
func (q numEval) Snapshot(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := rt.GetNumber(run, q.eval); e != nil {
		err = e
	} else {
		ret = generic.NewFloat(v)
	}
	return
}

func newNumValue(v interface{}) (ret snapper, err error) {
	switch a := v.(type) {
	case nil: // zero value for unhandled defaults in sqlite
		ret = staticValue{generic.Zero}
	case int:
		ret = staticValue{generic.NewInt(a)}
	case int64:
		ret = staticValue{generic.NewFloat(float64(a))}
	case float64:
		ret = staticValue{generic.NewFloat(a)}
	case []byte:
		var eval rt.NumberEval
		if e := bytesToEval(a, &eval); e != nil {
			err = e
		} else {
			ret = numEval{eval}
		}
	default:
		err = errutil.New("expected num value, got %v(%T)", v, v)
	}
	return
}
