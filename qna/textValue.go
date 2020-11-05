package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type textEval struct{ evalValue }

// GetText, or error if the underlying value isn't represented by a string.
func (q *textEval) GetText() (string, error) {
	return rt.GetText(q.run, q.eval.(rt.TextEval))
}

func newTextValue(run rt.Runtime, v interface{}) (ret rt.Value, err error) {
	switch a := v.(type) {
	case nil: // zero value for unhandled defaults in sqlite
		ret = generic.NewString("")
	case string:
		ret = generic.NewString(a)
	case []byte:
		var eval rt.TextEval
		if e := bytesToEval(a, &eval); e != nil {
			err = e
		} else {
			ret = &textEval{evalValue{run: run, eval: eval}}
		}
	default:
		err = errutil.New("expected text value, got %v(%T)", v, v)
	}
	return
}
