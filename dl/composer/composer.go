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
	OpenStrings             bool     // for str types, whether any value is permitted
	Strings                 []string // values for str types, generates tokens, labels, and selectors.
	Stub                    bool     // generate a custom loading struct.
	Locals                  []string
	Fluent                  *Fluid
}

func (x *Spec) UsesStr() bool {
	return x.OpenStrings || len(x.Strings) > 0
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
	// a non-top level command, basically any eval
	Function
	// commands that fill interfaces only in a small set of cases
	// example. the elseIf in an if.
	Selector
)

// Fluid provide extra info for displaying fluent commands.
type Fluid struct {
	Name string // if empty, use the type name
	Role Role
}
