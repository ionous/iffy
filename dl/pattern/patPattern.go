package pattern

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
)

type Pattern struct {
	Name    string
	Params  []term.Preparer
	Locals  []term.Preparer
	Returns term.Preparer
	Rules   []*Rule
}

// setup the runtime parameter info with our stored parameter info
func (ps *Pattern) ComputeParams(run rt.Runtime, terms *term.Terms) (ret int, err error) {
	return prepareList(run, ps.Params, terms)
}

func (ps *Pattern) ComputeLocals(run rt.Runtime, terms *term.Terms) (ret int, err error) {
	return prepareList(run, ps.Locals, terms)
}

func (ps *Pattern) ComputeReturn(run rt.Runtime, terms *term.Terms) (ret string, err error) {
	if res := ps.Returns; res != nil {
		if e := res.Prepare(run, terms); e != nil {
			err = e
		} else {
			ret = res.String()
		}
	}
	return
}

func (ps *Pattern) Execute(run rt.Runtime) (err error) {
	if inds, e := splitRules(run, ps.Rules); e != nil {
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

func prepareList(run rt.Runtime, list []term.Preparer, terms *term.Terms) (ret int, err error) {
	for _, n := range list {
		if e := n.Prepare(run, terms); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret++
		}
	}
	return
}
