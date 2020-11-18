package test

import (
	"github.com/ionous/errutil"
	g "github.com/ionous/iffy/rt/generic"
)

func SetRecord(d *g.Record, pairs ...interface{}) (err error) {
	for i, cnt := 0, len(pairs); i < cnt; i = i + 2 {
		if n, ok := pairs[0].(string); !ok {
			err = errutil.New("couldnt convert field")
		} else if v, e := g.ValueOf(pairs[1]); e != nil {
			err = errutil.New("couldnt convert value", e)
			break
		} else if e := d.SetNamedField(n, v); e != nil {
			err = e
			break
		}
	}
	return
}
