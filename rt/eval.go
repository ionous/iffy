package rt

import (
	"io"

	"github.com/ionous/errutil"
)

const StreamEnd errutil.Error = "stream end"
const StreamExceeded errutil.Error = "stream exceeded"

// Execute runs a bit of code that has no return value.
type Execute interface {
	// fix: rename to Run() to simplify look of Execute.Execute with embedded
	Execute(Runtime) error
}

// TextWriter runs a bit of code that writes into w.
type TextWriter interface {
	WriteText(Runtime, io.Writer) error
}

// BoolEval represents some boolean logic expression.
type BoolEval interface {
	GetBool(Runtime) (bool, error)
}

// NumberEval represents some numeric expression.
type NumberEval interface {
	GetNumber(Runtime) (float64, error)
}

// ObjectEval identifies some particular object in the game by its unique name.
type ObjectEval interface {
	GetObject(Runtime) (string, error)
}

// NumListEval returns or generates a series of numbers.
type NumListEval interface {
	GetNumberStream(Runtime) (NumberStream, error)
}

// NumListEval returns or generates a series of strings.
type TextListEval interface {
	GetTextStream(Runtime) (TextStream, error)
}

// NumListEval identifies a series of objects in the game.
type ObjListEval interface {
	GetObjectStream(Runtime) (ObjectStream, error)
}

// NumberStream provides a way to iterate over a set of numbers.
type NumberStream interface {
	HasNext() bool
	GetNumber() (float64, error)
}

// TextStream provides a way to iterate over a set of strings.
type TextStream interface {
	HasNext() bool
	GetText() (string, error)
}

// ObjectStream provides a way to iterate over a set of object.
type ObjectStream interface {
	HasNext() bool
	GetObject() (string, error)
}
