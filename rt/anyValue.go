package rt

type Value interface {
	// GetBool, or error if the underlying value isn't a bool
	GetBool(Runtime) (bool, error)
	// GetNumber, or error if the underlying value isn't a number
	GetNumber(Runtime) (float64, error)
	// GetText, or error if the underlying value isn't represented by a string.
	GetText(Runtime) (string, error)
	// GetLen returns the number of elements in the underlying value if it's a slice,
	// otherwise this returns an error.
	GetLen(Runtime) (int, error)
	// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
	// otherwise this returns an error.
	GetIndex(Runtime, int) (Value, error)
	// GetFieldByName for values representing objects ( nouns or records ), errors otherwise.
	GetFieldByName(Runtime, string) (Value, error)
}
