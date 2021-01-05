package generic

import (
	"github.com/ionous/iffy/affine"
)

// Value represents any one of the built in types.
// It's similar to reflect.Value in golang's standard library.
type Value interface {
	// Affinity identifies the general category
	Affinity() affine.Affinity
	// Type name of the specific underlying
	Type() string
	// Bool, or panic if the underlying value isn't a bool
	Bool() bool
	// Float, or panic if the underlying value isn't a number
	Float() float64
	// Float, or panic if the underlying value isn't a number
	Int() int
	// String, or panic if the underlying value isn't represented by a string.
	String() string
	// Record, returns the underlying record structure for record values
	Record() *Record
	// FloatSlice, or panic if the underlying value isn't represented by a slice of floats.
	Floats() []float64
	// StringSlice, or panic if the underlying value isn't represented by a slice of strings.
	Strings() []string
	// RecordSlice, or panic if the underlying value isn't represented by a slice of values.
	// ( every value in the returned list should be a record of this value's Type() )
	Records() []*Record
	// Index returns the nth element of the underlying slice, where 0 is the first value;
	// otherwise this returns an panic.
	Index(int) Value
	// Len returns the number of elements in the underylying value if it's a slice,
	// otherwise this returns an panic.
	Len() int
	// FieldByName for values representing objects, errors otherwise ( or if the field doesnt exit ).
	FieldByName(string) (Value, error)
	// SetFieldByName to write values back into objects.
	// errors if the affinities dont match. if the field doesnt exist, or if the value doesnt represent an object.
	SetFieldByName(string, Value) error
	// panics if out of range or if the values are mismatched.
	SetIndex(int, Value)
	// Append adds a value or value list if the underlying value is a slice.
	// In golang, this is a package level function, presumably to mirror the built-in append()
	Append(Value)
	// Slice returns a new list containing the first index up to (not including) the second index
	// Unlike go, the slices are distinct.
	Slice(i, j int) (Value, error)
	// Splice cuts elements from start to end, adding new elements at the start of the cut point
	// As a special case, passing a nil value to add only cuts elements.
	// Returns the cut elements, or an error if the start and end indices are bad;
	// panics if the element/s to add are of an incompatible type.
	Splice(start, end int, add Value) (Value, error)
}
