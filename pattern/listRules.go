package pattern

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

type NumListRules []NumListRule
type TextListRules []TextListRule
type ExecListRules []ExecuteRule

// ListRule for any rule which can respond with multiple results.
type ListRule struct {
	Filters Filters
	Flags
}

// NumListRule responds with a stream of numbers when its filters are satisfied.
// It works in conjunction with NumListRules.
type NumListRule struct {
	ListRule
	rt.NumListEval
}

// TextListRule responds with a stream of text when its filters are satisfied.
// It works in conjunction with TextListRules.
type TextListRule struct {
	ListRule
	rt.TextListEval
}

// ExecuteListRule triggers a series of statements when its filters are satisfied.
// It works in conjunction with ExecListRules.
type ExecuteRule struct {
	ListRule
	rt.Execute
}

func (r *ListRule) Apply(run rt.Runtime) (ret Flags, err error) {
	if ok, e := rt.GetAllTrue(run, r.Filters); e != nil {
		err = e
	} else if !ok {
		ret = -1
	} else {
		ret = r.Flags
	}
	return
}

// Apply returns flags if the filters passed, -1 if they did not, error on any error.
func (ps NumListRules) Apply(run rt.Runtime, i int) (ret Flags, err error) {
	return ps[i].Apply(run)
}

func (ps NumListRules) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if inds, e := splitRules(run, ps, len(ps)); e != nil {
		err = e
	} else {
		it := numIterator{run, ps, inds, 0}
		ret = stream.NewNumberChain(&it)
	}
	return
}

// Apply returns flags if the filters passed, -1 if they did not, error on any error.
func (ps TextListRules) Apply(run rt.Runtime, i int) (ret Flags, err error) {
	return ps[i].Apply(run)
}

func (ps TextListRules) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if inds, e := splitRules(run, ps, len(ps)); e != nil {
		err = e
	} else {
		it := textIterator{run, ps, inds, 0}
		ret = stream.NewTextChain(&it)
	}
	return
}

// Apply returns flags if the filters passed, -1 if they did not, error on any error.
func (ps ExecListRules) Apply(run rt.Runtime, i int) (ret Flags, err error) {
	return ps[i].Apply(run)
}

func (ps ExecListRules) Execute(run rt.Runtime) (ret bool, err error) {
	if inds, e := splitRules(run, ps, len(ps)); e != nil {
		err = e
	} else {
		for _, i := range inds {
			exec := ps[i]
			if e := (exec.Execute).Execute(run); e != nil {
				err = e
				break
			}
			ret = true // any executed
		}
	}
	return
}
