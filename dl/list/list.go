package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	g "github.com/ionous/iffy/rt/generic"
)

var Slats = []composer.Slat{
	(*At)(nil),
	(*Each)(nil),
	(*Len)(nil),
	(*Map)(nil),
	(*Pop)(nil),
	(*Push)(nil),
	(*Set)(nil),
	(*Slice)(nil),
	(*Sort)(nil),
	(*Splice)(nil),
}

func cmdError(op composer.Slat, e error) error {
	return errutil.Append(&core.CommandError{Cmd: op}, e)
}

// can a be inserted into b
func IsInsertable(el, els g.Value) (okay bool) {
	okay = true // provisionally
	listAff := els.Affinity()
	if needAff := affine.Element(listAff); len(needAff) == 0 {
		okay = false
	} else if haveAff := el.Affinity(); haveAff != needAff {
		okay = false
	} else if haveAff == affine.Record {
		if elt, elst := el.Type(), els.Type(); elt != elst {
			okay = false
		}
	}
	return
}

type insertError struct {
	el, els g.Value
}

func (e insertError) Error() string {
	return errutil.Sprintf("%s of %q isn't insertable into %s of %q",
		e.el.Affinity(), e.el.Type(),
		e.els.Affinity(), e.els.Type())
}
