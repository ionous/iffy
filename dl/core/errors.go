package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
)

type CommandError struct {
	Cmd composer.Slat
}

func (e *CommandError) Error() string {
	cmd := e.Cmd.Compose()
	return errutil.Sprintf("error in command %q", cmd.Name)
}

func cmdError(op composer.Slat, e error) error {
	return errutil.Append(&CommandError{op}, e)
}

type UnknownObject string

func (e UnknownObject) Error() string {
	return errutil.Sprintf("Unknown object %q", string(e))
}
