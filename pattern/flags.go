package pattern

// Flags controlling how individual list rules ( which each respond with multiple results ) work together.
type Flags int

//go:generate stringer -type=Flags
const (
	Infix Flags = iota
	Prefix
	Postfix
)
