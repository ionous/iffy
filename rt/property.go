package rt

// PropertyType enumerates the available kinds of fields on an object.
//go:generate stringer -type=PropertyType
type PropertyType int

// Enumeration of all possible object fields.
// Array is a bit field marker, the ordering of the constants is tuned to allow stringer to produce good results for array or'd fields
const (
	InvalidProperty PropertyType = iota
	Number                       // float64
	Text                         // string
	Pointer                      // *Object
	Array                        // array is a bit field marker
	NumberArray                  // []float64 NumberArray= Number|Array
	TextArray                    // []string
	PointerArray                 // []*Object
	State                        // int; note: states dont have arrays
)

// Property provides information on the fields of an object.
type Property interface {
	GetId() string
	GetType() PropertyType
}
