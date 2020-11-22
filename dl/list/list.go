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

// can el be inserted into els?
func IsInsertable(el, els g.Value) (okay bool) {
	return isInsertable(els, el.Affinity(), el.Type())
}

// can el be appended to els?
// this is similar to IsInsertable, except that the el can itself be a list.
func IsAppendable(el, els g.Value) (okay bool) {
	elAff := el.Affinity()
	if unlist := affine.Element(elAff); len(unlist) > 0 {
		elAff = unlist
	}
	return isInsertable(els, elAff, el.Type())
}

func isInsertable(els g.Value, haveAff affine.Affinity, haveType string) (okay bool) {
	okay = true // provisionally
	listAff := els.Affinity()
	if needAff := affine.Element(listAff); len(needAff) == 0 {
		okay = false // els was not actually a list
	} else if haveAff != needAff {
		okay = false // the element affinities dont match
	} else if haveAff == affine.Record && haveType != els.Type() {
		okay = false // the record types dont match
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
