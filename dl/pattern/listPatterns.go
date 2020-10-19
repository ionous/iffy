package pattern

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/chain"
)

type NumListPattern struct {
	CommonPattern
	Rules []*NumListRule
}
type TextListPattern struct {
	CommonPattern
	Rules []*TextListRule
}
type ActivityPattern struct {
	CommonPattern
	Rules []*ExecuteRule
}

// ApplyByIndex returns flags if the filters passed, -1 if they did not, error on any error.
func (ps *NumListPattern) ApplyByIndex(run rt.Runtime, i int) (ret Flags, err error) {
	return ps.Rules[i].GetFlags(run)
}

func (ps *NumListPattern) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if inds, e := splitNumbers(run, ps.Rules); e != nil {
		err = e
	} else {
		// fix: simplify this, doesnt need the iterators
		it := numIterator{run, ps, inds, 0}
		x := chain.NewStreamOfStreams(&it)
		ret, err = rt.CompactNumbers(x, nil)
	}
	return
}

// ApplyByIndex returns flags if the filters passed, -1 if they did not, error on any error.
func (ps *TextListPattern) ApplyByIndex(run rt.Runtime, i int) (ret Flags, err error) {
	return ps.Rules[i].GetFlags(run)
}

func (ps *TextListPattern) GetTextList(run rt.Runtime) (ret []string, err error) {
	if inds, e := splitText(run, ps.Rules); e != nil {
		err = e
	} else {
		it := textIterator{run, ps, inds, 0}
		// fix: simplify this, doesnt need the iterators
		x := chain.NewStreamOfStreams(&it)
		ret, err = rt.CompactTexts(x, nil)
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
