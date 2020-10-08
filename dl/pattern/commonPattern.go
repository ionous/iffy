package pattern

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type Pattern interface {
	// setup the runtime parameter info with our stored parameter info
	Prepare(rt.Runtime, *Parameters) error
	ComputeLocals(rt.Runtime, *Parameters) error
	GetParameterName(int) (string, error)
}

// fix: the duplication of this and the name, prologue parameters indicates that
// the structure is inverted -- there should probably be one common pattern struct
// with a rules interface implemented by lists of Text, etc rules.
type CommonPattern struct {
	Name     string
	Prologue []Parameter
	Locals   []Parameter
}

// setup the runtime parameter info with our stored parameter inf
func (ps *CommonPattern) Prepare(run rt.Runtime, parms *Parameters) (err error) {
	return prepareList(run, ps.Prologue, parms)
}

func (ps *CommonPattern) ComputeLocals(run rt.Runtime, parms *Parameters) (err error) {
	return prepareList(run, ps.Locals, parms)
}

func prepareList(run rt.Runtime, list []Parameter, parms *Parameters) (err error) {
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
