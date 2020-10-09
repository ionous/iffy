package pattern

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
)

type Pattern interface {
	// setup the runtime parameter info with our stored parameter info
	Prepare(rt.Runtime, *term.Terms) error
	ComputeLocals(rt.Runtime, *term.Terms) error
	GetParameterName(int) (string, error)
}

// fix: the duplication of this and the name, prologue parameters indicates that
// the structure is inverted -- there should probably be one common pattern struct
// with a rules interface implemented by lists of Text, etc rules.
type CommonPattern struct {
	Name     string
	Prologue []term.Preparer
	Locals   []term.Preparer
}

// setup the runtime parameter info with our stored parameter inf
func (ps *CommonPattern) Prepare(run rt.Runtime, parms *term.Terms) (err error) {
	return prepareList(run, ps.Prologue, parms)
}

func (ps *CommonPattern) ComputeLocals(run rt.Runtime, parms *term.Terms) (err error) {
	return prepareList(run, ps.Locals, parms)
}

func prepareList(run rt.Runtime, list []term.Preparer, parms *term.Terms) (err error) {
	for _, n := range list {
		if e := n.Prepare(run, parms); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (ps *CommonPattern) GetParameterName(idx int) (ret string, err error) {
	if idx < 0 || idx >= len(ps.Prologue) {
		err = errutil.New("indexed parameter out of range", idx)
	} else {
		// alt: we could use the database to search GetFieldByIndex
		p := ps.Prologue[idx]
		// preliminarily, the parameters are just their names.
		ret = p.String()
	}
	return
}
