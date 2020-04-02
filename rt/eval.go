package rt

import "github.com/ionous/errutil"

const StreamExceeded errutil.Error = "stream exceeded"

// Execute runs a bit of code that has no return value.
type Execute interface {
	// fix: rename to Run() to simplify look of Execute.Execute with embedded
	Execute(Runtime) error
}

// BoolEval represents some boolean logic expression.
type BoolEval interface {
	GetBool(Runtime) (bool, error)
}

// NumberEval represents some numeric expression.
type NumberEval interface {
	GetNumber(Runtime) (float64, error)
}

// TextEval runs a bit of code that writes into w.
type TextEval interface {
	GetText(Runtime) (string, error)
}

// NumListEval returns or generates a series of numbers.
type NumListEval interface {
	GetNumberStream(Runtime) (NumberStream, error)
}

// NumListEval returns or generates a series of strings.
type TextListEval interface {
	GetTextStream(Runtime) (TextStream, error)
}

// NumberStream provides a way to iterate over a set of numbers.
type NumberStream interface {
	// HasNext returns true if the iterator can be safely advanced.
	HasNext() bool
	// GetNumber advances the iterator.
	GetNumber() (float64, error)
}

// TextStream provides a way to iterate over a set of strings.
type TextStream interface {
	// HasNext returns true if the iterator can be safely advanced.
	HasNext() bool
	// GetText advances the iterator.
	GetText() (string, error)
}

// StreamCount provides an optional interface for determining the number of elements in a stream.
type StreamCount interface {
	// Count returns the remaining length of the stream.
	Count() int
}
