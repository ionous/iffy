package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type numEval struct{ evalValue }

// GetText, or error if the underlying value isn't represented by a string.
func (q *numEval) GetNumber() (float64, error) {
	return rt.GetNumber(q.run, q.eval.(rt.NumberEval))
}

func newNumValue(run rt.Runtime, v interface{}) (ret rt.Value, err error) {
	switch a := v.(type) {
	case nil: // zero value for unhandled defaults in sqlite
		ret = &generic.Float{}
	case int:
		ret = &generic.Int{Value: a}
	case int64:
		ret = &generic.Float{Value: float64(a)}
	case float64:
		ret = &generic.Float{Value: a}
	case []byte:
		var eval rt.NumberEval
		if e := bytesToEval(a, &eval); e != nil {
			err = e
		} else {
			ret = &numEval{evalValue{run: run, eval: eval}}
		}
	default:
		err = errutil.New("expected num value, got %v(%T)", v, v)
	}
	return
}
