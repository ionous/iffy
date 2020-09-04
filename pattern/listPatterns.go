package pattern

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

type NumListPattern struct {
	Name  string
	Rules []*NumListRule
}
type TextListPattern struct {
	Name  string
	Rules []*TextListRule
}
type ActivityPattern struct {
	Name  string
	Rules []*ExecuteRule
}

// ApplyByIndex returns flags if the filters passed, -1 if they did not, error on any error.
func (ps *NumListPattern) ApplyByIndex(run rt.Runtime, i int) (ret Flags, err error) {
	return ps.Rules[i].GetFlags(run)
}

func (ps *NumListPattern) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if inds, e := splitNumbers(run, ps.Rules); e != nil {
		err = e
	} else {
		it := numIterator{run, ps, inds, 0}
		ret = stream.NewNumberChain(&it)
	}
	return
}

// ApplyByIndex returns flags if the filters passed, -1 if they did not, error on any error.
func (ps *TextListPattern) ApplyByIndex(run rt.Runtime, i int) (ret Flags, err error) {
	return ps.Rules[i].GetFlags(run)
}

func (ps *TextListPattern) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	if inds, e := splitText(run, ps.Rules); e != nil {
		err = e
	} else {
		it := textIterator{run, ps, inds, 0}
		ret = stream.NewTextChain(&it)
	}
	return
}

// ApplyByIndex returns flags if the filters passed, -1 if they did not, error on any error.
func (ps *ActivityPattern) ApplyByIndex(run rt.Runtime, i int) (ret Flags, err error) {
	return ps.Rules[i].GetFlags(run)
}

func (ps *ActivityPattern) Execute(run rt.Runtime) (err error) {
	if inds, e := splitExe(run, ps.Rules); e != nil {
		err = e
	} else {
		for _, i := range inds {
			if e := rt.RunOne(run, ps.Rules[i].Execute); e != nil {
				err = e
				break
			}
			// NOTE: if we need to differentiate between "ran" and "not found",
			// "didnt run" should probably become an error code.
		}
	}
	return
}
