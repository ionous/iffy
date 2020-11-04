package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

var Slats = []composer.Slat{
	(*At)(nil),
	(*Each)(nil),
	(*Len)(nil),
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
