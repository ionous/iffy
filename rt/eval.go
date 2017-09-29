package rt

import (
	"github.com/ionous/errutil"
)

const StreamEnd errutil.Error = "stream end"
const StreamExceeded errutil.Error = "stream exceeded"

type Execute interface {
	// fix: rename to Run() to simplify look of Execute.Execute with embedded
	Execute(Runtime) error
}
type BoolEval interface {
	GetBool(Runtime) (bool, error)
}
type NumberEval interface {
	GetNumber(Runtime) (float64, error)
}
type TextEval interface {
	GetText(Runtime) (string, error)
}
type ObjectEval interface {
	GetObject(Runtime) (Object, error)
}
type NumListEval interface {
	GetNumberStream(Runtime) (NumberStream, error)
}
type TextListEval interface {
	GetTextStream(Runtime) (TextStream, error)
}
type ObjListEval interface {
	GetObjectStream(Runtime) (ObjectStream, error)
}

type NumberStream interface {
	HasNext() bool
	GetNumber() (float64, error)
}
type TextStream interface {
	HasNext() bool
	GetText() (string, error)
}
type ObjectStream interface {
	HasNext() bool
	GetObject() (Object, error)
}
