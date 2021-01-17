package pattern

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/chain"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type xNumListPattern struct {
	CommonPattern
	Rules []*xNumListRule
}
type xTextListPattern struct {
	CommonPattern
	Rules []*xTextListRule
}
type ActivityPattern struct {
	CommonPattern
	Rules []*ExecuteRule
}

func (ps *xNumListPattern) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if inds, e := xsplitNumbers(run, ps.Rules); e != nil {
		err = e
	} else {
		// fix: simplify this, doesnt need the iterators
		it := xnumIterator{run, ps, inds, 0}
		x := chain.NewStreamOfStreams(&it)
		ret, err = g.CompactNumbers(x, nil)
	}
	return
}

func (ps *xTextListPattern) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if inds, e := xsplitText(run, ps.Rules); e != nil {
		err = e
	} else {
		it := xtextIterator{run, ps, inds, 0}
		// fix: simplify this, doesnt need the iterators
		x := chain.NewStreamOfStreams(&it)
		ret, err = g.CompactTexts(x, nil)
	}
	return
}

func (ps *ActivityPattern) Execute(run rt.Runtime) (err error) {
	if inds, e := splitExe(run, ps.Rules); e != nil {
		err = e
	} else {
		for _, i := range inds {
			if e := safe.Run(run, ps.Rules[i].Execute); e != nil {
				err = e
				break
			}
			// NOTE: if we need to differentiate between "ran" and "not found",
			// "didnt run" should probably become an error code.
		}
	}
	return
}
