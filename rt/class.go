package rt

// Class describes a shared class of Objects.
type Class interface {
	// GetId returns the unique identifier for this type.
	GetId() string
	// GetName returns the user defined name for this type.
	GetName() string
	// GetParentType returns false for types if no parent;
	GetParent() (Class, bool)
	// IsCompatible checks whether this class can be used as the named class.
	// ie. this is the named class, or a descendant of it.
	IsCompatible(string) bool
}
