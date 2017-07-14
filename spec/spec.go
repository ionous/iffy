package spec

// Factory creates new command writers.
type Factory interface {
	NewSpec(name string) (Spec, error)
	NewSpecs() (Specs, error)
}

// Spec writes to the fields of a command.
type Spec interface {
	// Position adds a new positional argument.
	// Positional arguments are guaranteed to precede keyword arguments.
	Position(arg interface{}) error
	// Assign adds a new keyword argument.
	// Keyword are guaranteed to follow any positional arguments.
	Assign(key string, value interface{}) error
}

// Specs writes to an array of commands.
type Specs interface {
	AddElement(Spec) error
}

// // Slot helps users write to blocks.
// type Slot interface {
// 	Cmd(name string, args ...interface{}) Block
// 	// Cmds starts a new array of commands.
// 	// It takes initial members of that block.
// 	// The block must eventually be terminated with End().
// 	Cmds(cmds ...Block) Block

// 	// Val specifies a single literal value: whether one primitive value or one array of primitive values. It retuns the current block
// 	Val(val interface{}) Block
// }

// // Block helps users build trees of commands.
// type Block interface {
// 	// Begin starts a new parameter block. Usually used as:
// 	//  if c.Cmd("name").Begin() {
// 	//    c.End()
// 	//  }
// 	Begin() (okay bool)

// 	// Slot writes to the next slot in this block
// 	Slot
// 	Param(field string) Slot

// 	// End terminates a block.
// 	// End must be called block must eventually be terminated with End().
// 	End()
// }
