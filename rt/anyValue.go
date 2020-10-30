package rt

import "github.com/ionous/iffy/affine"

// Value represents any one of the built in types.
// While a raw interface{} would work as well
// this helps internalize and standardize error codes.
type Value interface {
	// Affinity identifies the general category
	Affinity() affine.Affinity
	// Type name of the specific underlying
	Type() string
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
	// GetField for values representing objects, errors otherwise ( or if the field doesnt exit ).
	GetField(string) (Value, error)
	// SetField to write values back into objects,
	// errors if the affinities dont match. if the field doesnt exist, or if the value doesnt represent an object.
	SetField(string, Value) error
}
