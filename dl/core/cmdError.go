package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
)

type CmdError struct {
	Cmd composer.Slat
	Err error
}

func (e CmdError) Error() string {
	cmd := e.Cmd.Compose()
	return errutil.Sprintf("%s encountered %v", cmd.Name, e.Err)
}
