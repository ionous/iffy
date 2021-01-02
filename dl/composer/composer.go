package composer

// Slot definition for display in the composer.
type Slot struct {
	Name  string
	Type  interface{} // nil instance, ex. (*core.Comparator)(nil)
	Desc  string
	Group string // display group(s)
}

// Spec definition for display in composer.
type Spec struct {
	Name, Spec, Group, Desc string
	Choices                 []string // for str types
	Stub                    bool     // generate a custom loading struct.
	Locals                  []string
	Fluent                  *Fluid
}

type Composer interface {
	Compose() Spec
}

// for highlighting info
// while this could be determined algorithmically, it would be a bunch of extra code
type Role int

//go:generate stringer -type=Role
const (
	// a top-level function, basically execute
	Command Role = iota + 1
	// any sub function, basically any eval
	Function
	// keyword choices
	Selector
)

// Fluid provide extra info for displaying fluent commands.
type Fluid struct {
	Name string // if empty, use the type name
	Role Role
}
