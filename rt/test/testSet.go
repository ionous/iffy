package test

import (
	"github.com/ionous/errutil"
	g "github.com/ionous/iffy/rt/generic"
)

func SetRecord(d *g.Record, pairs ...interface{}) (err error) {
	for i, cnt := 0, len(pairs); i < cnt; i = i + 2 {
		if n, ok := pairs[0].(string); !ok {
			err = errutil.New("couldnt convert field")
		} else {

			if v, e := ValueOf(pairs[1]); e != nil {
				err = errutil.New("couldnt convert value", e)
				break
			} else if e := d.SetNamedField(n, v); e != nil {
				err = e
				break
			}
		}
	}
	return
}

// ValueOf returns a new generic value wrapper based on analyzing the passed value.
func ValueOf(i interface{}) (ret g.Value, err error) {
	switch i := i.(type) {
	case bool:
		ret = g.BoolOf(i)
	case int:
		ret = g.IntOf(i)
	case float64:
		ret = g.FloatOf(i)
	case string:
		ret = g.StringOf(i)
	case []float64:
		ret = g.FloatsOf(i)
	case []string:
		ret = g.StringsOf(i)
	case *g.Record:
		ret = g.RecordOf(i)
	default:
		err = errutil.New("unhandled value %v(%T)", i, i)
	}
	return
}
