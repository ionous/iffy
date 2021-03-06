package debug

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

var Slats = []composer.Composer{
	(*DoNothing)(nil),
	(*Log)(nil),
	(*Level)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}
