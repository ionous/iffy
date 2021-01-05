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

func cmdError(op composer.Composer, e error) error {
	return errutil.Append(&CommandError{Cmd: op}, e)
}

func cmdErrorCtx(op composer.Composer, ctx string, e error) error {
	return errutil.Append(&CommandError{Cmd: op, Ctx: ctx}, e)
}
