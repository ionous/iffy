package pattern

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

// ListRule for any rule which can respond with multiple results.
type ListRule struct {
	Filter rt.BoolEval
	Flags
}

// NumListRule responds with a stream of numbers when its filters are satisfied.
type NumListRule struct {
	ListRule
	rt.NumListEval
}

// TextListRule responds with a stream of text when its filters are satisfied.
type TextListRule struct {
	ListRule
	rt.TextListEval
}

// ExecuteListRule triggers a series of statements when its filters are satisfied.
type ExecuteRule struct {
	ListRule
	rt.Execute
}

func (r *ListRule) GetFlags(run rt.Runtime) (ret Flags, err error) {
	if ok, e := safe.GetOptionalBool(run, r.Filter, true); e != nil {
		err = e
	} else if !ok.Bool() {
		ret = -1
	} else {
		ret = r.Flags
	}
	return
}
