package rt

type Value interface {
	// GetBool, or error if the underlying value isn't a bool
	GetBool(Runtime) (bool, error)
	// GetNumber, or error if the underlying value isn't a number
	GetNumber(Runtime) (float64, error)
	// GetText, or error if the underlying value isn't represented by a string.
	GetText(Runtime) (string, error)
	// GetNumber, or error if the underlying value isn't a number
	GetNumList(Runtime) ([]float64, error)
	// GetText, or error if the underlying value isn't represented by a string.
	GetTextList(Runtime) ([]string, error)
	// GetFieldByName for values representing objects ( nouns or records ), errors otherwise.
	GetFieldByName(Runtime, string) (Value, error)
}
