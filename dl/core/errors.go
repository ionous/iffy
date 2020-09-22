package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
)

type CommandError struct {
	Cmd composer.Slat
	Ctx string
}

func (e *CommandError) Error() string {
	cmd := e.Cmd.Compose()
	var padding rune
	if len(e.Ctx) > 0 {
		padding = ' '
	}
	return errutil.Sprintf("error in command %q%v%s", cmd.Name, padding, e.Ctx)
}

func cmdError(op composer.Slat, e error) error {
	return errutil.Append(&CommandError{Cmd: op}, e)
}

func cmdErrorCtx(op composer.Slat, Ctx string, e error) error {
	return errutil.Append(&CommandError{Cmd: op, Ctx: Ctx}, e)
}

type UnknownObject string

func (e UnknownObject) Error() string {
	return errutil.Sprintf("Unknown object %q", string(e))
}
