package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/reader"
)

type OpError struct {
	Op  composer.Slat
	At  reader.Position
	Err error
}

const UnhandledSwap = errutil.Error("unhandled swap")
const MissingSlot = errutil.Error("missing slot")
const InvalidValue = errutil.Error("invalid value")

func ImportError(op composer.Slat, at reader.Position, e error) error {
	return &OpError{op, at, e}
}

func (e *OpError) Error() string {
	op := e.Op.Compose()
	return errutil.Sprintf("%s in %s at %s", e.Err, op.Name, e.At.String())
}

func (e *OpError) Unwrap() error {
	return e.Err
}
