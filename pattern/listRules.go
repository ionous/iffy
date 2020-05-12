package pattern

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

type NumListRules []NumListRule
type TextListRules []TextListRule
type ExecRules []ExecuteRule

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
	Go []rt.Execute
}

func (*NumListRule) RuleDesc() RuleDesc {
	return RuleDesc{
		"num_list_rule",
		(*NumListRules)(nil),
	}
}

func (*TextListRule) RuleDesc() RuleDesc {
	return RuleDesc{
		"text_list_rule",
		(*TextListRules)(nil),
	}
}

func (*ExecuteRule) RuleDesc() RuleDesc {
	return RuleDesc{
		"execute_rule",
		(*ExecRules)(nil),
	}
}

func (r *ListRule) ApplyByIndex(run rt.Runtime) (ret Flags, err error) {
	if ok, e := rt.GetOptionalBool(run, r.Filter, true); e != nil {
		err = e
	} else if !ok {
		ret = -1
	} else {
		ret = r.Flags
	}
	return
}

// ApplyByIndex returns flags if the filters passed, -1 if they did not, error on any error.
func (ps NumListRules) ApplyByIndex(run rt.Runtime, i int) (ret Flags, err error) {
	return ps[i].ApplyByIndex(run)
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

// ApplyByIndex returns flags if the filters passed, -1 if they did not, error on any error.
func (ps TextListRules) ApplyByIndex(run rt.Runtime, i int) (ret Flags, err error) {
	return ps[i].ApplyByIndex(run)
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

func (ps ExecuteRule) Execute(run rt.Runtime) error {
	return rt.RunAll(run, ps.Go)
}

// ApplyByIndex returns flags if the filters passed, -1 if they did not, error on any error.
func (ps ExecRules) ApplyByIndex(run rt.Runtime, i int) (ret Flags, err error) {
	return ps[i].ApplyByIndex(run)
}

func (ps ExecRules) Execute(run rt.Runtime) (ret bool, err error) {
	if inds, e := splitRules(run, ps, len(ps)); e != nil {
		err = e
	} else {
		for _, i := range inds {
			if e := rt.RunAll(run, ps[i].Go); e != nil {
				err = e
				break
			}
			ret = true // any executed
		}
	}
	return
}
