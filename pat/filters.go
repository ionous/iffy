package pat

import (
	"github.com/ionous/iffy/rt"
)

type Filters []rt.BoolEval

func (fs Filters) GetBool(run rt.Runtime) (okay bool, err error) {
	// walk filters in the same direction that they are specified in their AllTrue list.
	i, cnt := 0, len(fs)
	for ; i < cnt; i++ {
		f := fs[i]
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
