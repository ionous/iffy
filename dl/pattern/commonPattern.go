package pattern

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
)

type Pattern interface {
	// setup the runtime parameter info with our stored parameter info
	ComputeParams(rt.Runtime, *term.Terms) (int, error)
	ComputeLocals(rt.Runtime, *term.Terms) (int, error)
	ComputeReturn(rt.Runtime, *term.Terms) (string, error)
}

// fix: the duplication of this and the name, prologue parameters indicates that
// the structure is inverted -- there should probably be one common pattern struct
// with a rules interface implemented by lists of Text, etc rules.
type CommonPattern struct {
	Name    string
	Params  []term.Preparer
	Locals  []term.Preparer
	Returns term.Preparer
}

// setup the runtime parameter info with our stored parameter info
func (ps *CommonPattern) ComputeParams(run rt.Runtime, terms *term.Terms) (ret int, err error) {
	return prepareList(run, ps.Params, terms)
}

func (ps *CommonPattern) ComputeLocals(run rt.Runtime, terms *term.Terms) (ret int, err error) {
	return prepareList(run, ps.Locals, terms)
}

func (ps *CommonPattern) ComputeReturn(run rt.Runtime, terms *term.Terms) (ret string, err error) {
	if res := ps.Returns; res != nil {
		if e := res.Prepare(run, terms); e != nil {
			err = e
		} else {
			ret = res.String()
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
