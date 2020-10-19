package rt

import "github.com/ionous/iffy/affine"

// Value represents any one of the built in types.
// While a raw interface{} would work as well
// this helps internalize and standardize error codes.
type Value interface {
	Affinity() affine.Affinity
	// GetBool, or error if the underlying value isn't a bool
	GetBool() (bool, error)
	// GetNumber, or error if the underlying value isn't a number
	GetNumber() (float64, error)
	// GetText, or error if the underlying value isn't represented by a string.
	GetText() (string, error)
	// GetNumber, or error if the underlying value isn't a number
	GetNumList() ([]float64, error)
	// GetText, or error if the underlying value isn't represented by a string.
	GetTextList() ([]string, error)
	// GetLen returns the number of elements in the underlying value if it's a slice,
	// otherwise this returns an error.
	GetLen() (int, error)
	// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
	// otherwise this returns an error.
	GetIndex(int) (Value, error)
	// GetFieldByName for values representing objects ( nouns or records ), errors otherwise.
	GetFieldByName(string) (Value, error)
}
