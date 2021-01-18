package pattern

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

// Rule triggers a series of statements when its filters are satisfied.
type Rule struct {
	Filter rt.BoolEval
	Flags
	rt.Execute
}

func (r *Rule) GetFlags(run rt.Runtime) (ret Flags, err error) {
	if ok, e := safe.GetOptionalBool(run, r.Filter, true); e != nil {
		err = e
	} else if !ok.Bool() {
		ret = -1
	} else {
		ret = r.Flags
	}
	return
}
