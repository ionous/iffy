package spec

// Spec writes to the fields of a command.
type Spec interface {
	// Position adds a new positional argument.
	// Positional arguments are guaranteed to precede keyword arguments.
	Position(arg interface{}) error
	// Assign adds a new keyword argument.
	// Keyword are guaranteed to follow any positional arguments.
	Assign(k string, v interface{}) error
}

// Specs writes to an array of commands.
type Specs interface {
	AddElement(Spec) error
}

// Slot helps users write to blocks.
type Slot interface {
	// Begin starts a new block.
	// If chained from a Cmd, this starts a block of command parameters,
	// Otherwise, it starts an block of commands.
	Begin() bool
	// Cmd starts a new command block.
	// It returns the *current+* block.
	// To gain access to the newly created block, chain this call with a call to Begin().
	Cmd(name string, args ...interface{}) Block
	// Val specifies a single literal value: whether one primitive value or one array of primitive values.
	// It returns the current block.
	Val(val interface{}) Block
}

// Block helps users build trees of commands.
type Block interface {
	// Slot writes to the next slot in the current block
	Slot
	// Param targets a specific field in the current block.
	Param(field string) Slot
	// End terminates a block.
	// End must be called block must eventually be terminated with End().
	End()
}
