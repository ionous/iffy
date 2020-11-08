package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type textEval struct{ eval rt.TextEval }

// GetText, or error if the underlying value isn't represented by a string.
func (q textEval) Snapshot(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := rt.GetText(run, q.eval); e != nil {
		err = e
	} else {
		ret = generic.NewString(v)
	}
	return
}

func newTextValue(v interface{}) (ret snapper, err error) {
	switch a := v.(type) {
	case nil: // zero value for unhandled defaults in sqlite
		ret = staticValue{generic.Empty}
	case string:
		ret = staticValue{generic.NewString(a)}
	case []byte:
		var eval rt.TextEval
		if e := bytesToEval(a, &eval); e != nil {
			err = e
		} else {
			ret = textEval{eval}
		}
	default:
		err = errutil.New("expected text value, got %v(%T)", v, v)
	}
	return
}
