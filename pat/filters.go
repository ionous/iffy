package pat

import (
	"github.com/ionous/iffy/rt"
)

type Filters []rt.BoolEval

// NOTE: filters are tested in reverse order.
func (fs Filters) GetBool(run rt.Runtime) (okay bool, err error) {
	i, cnt := 0, len(fs)
	for ; i < cnt; i++ {
		f := fs[cnt-i-1]
		if ok, e := f.GetBool(run); e != nil {
			err = e
			break
		} else if !ok {
			break
		}
	}
	okay = i == cnt
	return
}
