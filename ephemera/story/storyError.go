package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/reader"
)

type OpError struct {
	Op  composer.Composer
	At  reader.Position
	Err error
}

const UnhandledSwap = errutil.Error("unhandled swap")
const MissingSlot = errutil.Error("missing slot")
const InvalidValue = errutil.Error("invalid value")

func ImportError(op composer.Composer, at reader.Position, e error) error {
	return &OpError{op, at, e}
}

func (e *OpError) Error() string {
	return errutil.Sprintf("%s in %s at %s", e.Err, composer.SpecName(e.Op), e.At.String())
}

func (e *OpError) Unwrap() error {
	return e.Err
}
