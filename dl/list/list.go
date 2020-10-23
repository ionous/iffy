package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

var Slats = []composer.Slat{
	(*At)(nil),
	(*Len)(nil),
	(*Pop)(nil),
	(*Push)(nil),
	(*Slice)(nil),
	(*Splice)(nil),
}

func cmdError(op composer.Slat, e error) error {
	return errutil.Append(&core.CommandError{Cmd: op}, e)
}

// for each
// sort ( w/ pattern )
