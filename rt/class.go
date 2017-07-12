package rt

// Class describes a shared class of Objects.
type Class interface {
	// GetId returns the unique identifier for this type.
	GetId() string
	// GetName returns the user defined name for this type.
	GetName() string
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
