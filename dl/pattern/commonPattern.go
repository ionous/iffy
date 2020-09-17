package pattern

import "github.com/ionous/errutil"

type Pattern interface {
	Prepare(ps Parameters)
	GetParameterName(int) (string, error)
}

func (ps *CommonPattern) Prepare(parms Parameters) {
	for _, n := range ps.Prologue {
		n.Prepare(parms)
	}
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

// fix: the duplication of this and the name, prologue parameters indicates that
// the structure is inverted -- there should probably be one common pattern struct
// with a rules interface implemented by lists of Text, etc rules.
type CommonPattern struct {
	Name     string
	Prologue []Parameter
}
