package render

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

var Slats = []composer.Composer{
	(*RenderField)(nil),
	(*RenderName)(nil),
	(*RenderRef)(nil),
	(*RenderTemplate)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}
