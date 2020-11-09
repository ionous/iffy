package generic

// Iterator provides a way to iterate over a stream of values.
// The underlying implementation and values returned depends on the stream.
type Iterator interface {
	// HasNext returns true if the iterator can be safely advanced.
	HasNext() bool
	// GetNext returns the next value in the stream and advances the iterator.
	GetNext() (Value, error)
}
