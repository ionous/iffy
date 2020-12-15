package render

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

var Slats = []composer.Composer{
	(*Name)(nil),
	(*Template)(nil),
}

func cmdError(op composer.Composer, e error) error {
	return errutil.Append(&core.CommandError{Cmd: op}, e)
}
