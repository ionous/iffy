package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type boolEval struct{ eval rt.BoolEval }

func (q boolEval) Snapshot(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := rt.GetBool(run, q.eval); e != nil {
		err = e
	} else {
		ret = generic.NewBool(v)
	}
	return
}

func newBoolValue(v interface{}) (ret snapper, err error) {
	switch a := v.(type) {
	case nil: // zero value for unhandled defaults in sqlite
		ret = staticValue{generic.False}
	case bool:
		ret = staticValue{generic.NewBool(a)}
	case int64: // sqlite, boolean values can be represented as 1/0
		ret = staticValue{generic.NewBool(a != 0)}
	case []byte:
		var eval rt.BoolEval
		if e := bytesToEval(a, &eval); e != nil {
			err = e
		} else {
			ret = &boolEval{eval}
		}
	default:
		err = errutil.Fmt("expected boolean value, got %v(%T)", v, v)
	}
	return
}
