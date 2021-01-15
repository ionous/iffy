package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
)

type CommandError struct {
	Cmd composer.Composer
	Ctx string
}

func (e *CommandError) Error() string {
	name := composer.SpecName(e.Cmd)
	var padding string
	if len(e.Ctx) > 0 {
		padding = " "
	}
	return errutil.Sprintf("error in command %q%s%s", name, padding, e.Ctx)
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &CommandError{Cmd: op})
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	return errutil.Append(err, &CommandError{Cmd: op, Ctx: ctx})
}
