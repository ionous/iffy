package rt

import "github.com/ionous/iffy/affine"

// Value represents any one of the built in types.
// It's similar to reflect.Value in golang's standard library.
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
	// GetNumList, or error if the underlying value isn't represented by a slice of floats.
	GetNumList() ([]float64, error)
	// GetTextList, or error if the underlying value isn't represented by a slice of strings.
	GetTextList() ([]string, error)
	// GetRecordList, or error if the underlying value isn't represented by a slice of values.
	// ( every value in the returned list should be a record of this value's Type() )
	GetRecordList() ([]Value, error)
	// GetIndex returns the nth element of the underlying slice, where 0 is the first value;
	// otherwise this returns an error.
	GetIndex(int) (Value, error)
	// GetLen returns the number of elements in the underlying value if it's a slice,
	// otherwise this returns an error.
	GetLen() (int, error)
	// GetNamedField for values representing objects, errors otherwise ( or if the field doesnt exit ).
	GetNamedField(string) (Value, error)
	// SetNamedField to write values back into objects.
	// errors if the affinities dont match. if the field doesnt exist, or if the value doesnt represent an object.
	SetNamedField(string, Value) error
	//
	SetIndexedValue(int, Value) error
	// Append adds a value or value list if the underlying value is a slice.
	// In golang, this is a package level function, presumably to mirror the built-in append()
	Append(Value) (Value, error)
	// Resize modifies the length of a list value.
	Resize(int) (Value, error)
	// Slice returns a new list containing the first index up to (not including) the second index
	// As in golang, the two slices are "aliased" -- they point to the same memory.
	Slice(i, j int) (Value, error)
}
