package ref

// Model provides the starting point for all game objects and types.
type Model interface {
	NumClass() int
	ClassNum(int) Class
	GetClass(string) (Class, bool)

	NumObject() int
	ObjectNum(int) Object
	GetObject(string) (Object, bool)
}

// Class describes a shared class of Objects.
type Class interface {
	// GetId returns the unique identifier for this types.
	GetId() string
	// GetParentType returns false for types if no parent;
	GetParent() (Class, bool)
	// NumProperty returns the number of indexable properties.
	// The number of available properties for a given Class never changes at runtime.
	NumProperty() int
	// PropertyNum returns the indexed property.
	// Panics if the index is greater than Number.
	PropertyNum(int) Property
	// GetProperty by name.
	GetProperty(string) (Property, bool)
	// GetPropertyByChoice evaluates all properties to find an enumeration which can store the passed choice
	GetPropertyByChoice(string) (Property, bool)
	// IsCompatible checks whether this Class is a child of the passed named parent.
	IsCompatible(string) bool
}

// Object represents a tangible or intangible piece of the game world.
type Object interface {
	// GetId returns the unique identifier for this Object.
	GetId() string
	// GetClass returns the variety of object.
	GetClass() Class
	// GetValue stores the value into the pointer pv.
	// Values include meta.Objects for relations and pointers, numbers, and text. For numbers, pv can be any numberic type: float64, int, etc.
	GetValue(name string, pv interface{}) error
	// GetValue can return error when the value violates a property constraint,
	// if the value is not of the requested type, or if the targeted property holder is read-only. Read-only values include the "many" side of a relation.
	SetValue(name string, v interface{}) error
}

// Property provides information on the fields of an object.
type Property interface {
	GetId() string
	GetType() PropertyType
}

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
