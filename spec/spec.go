package spec

type Factory interface {
	NewSpec(name string) (Spec, error)
	NewSpecs() (Specs, error)
}

type Spec interface {
	// Position adds a new positional argument.
	// Positional arguments are guaranteed to precede keyword arguments.
	Position(arg interface{}) error
	// Assign adds a new keyword argument.
	// Keyword are guaranteed to follow any positional arguments.
	Assign(key string, value interface{}) error
}

type Specs interface {
	AddElement(Spec) error
}
