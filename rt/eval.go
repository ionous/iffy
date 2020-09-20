package rt

// Execute runs a bit of code that has no return value.
type Execute interface {
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
	GetNumberStream(Runtime) (Iterator, error)
}

// NumListEval returns or generates a series of strings.
type TextListEval interface {
	GetTextStream(Runtime) (Iterator, error)
}

// Iterator provides a way to iterate over a stream of values.
// The underlying implementation and values returned depends on the stream.
type Iterator interface {
	// HasNext returns true if the iterator can be safely advanced.
	HasNext() bool
	// GetNext returns the next value in the stream and advances the iterator.
	GetNext() (Value, error)
}

// StreamCount optionally implemented for iterators to determine the length of remaining stream.
type StreamCount interface {
	Remaining() int
}
