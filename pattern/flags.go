package pattern

// Flags control how patterns involving lists generate results.
// Rules which pass their filters get sorted and their output gets generated in the flagged order.
type Flags int

//go:generate stringer -type=Flags
const (
	Infix   Flags = iota // default, stops sorting all other rules
	Prefix               // sorts the rule execution to the font
	Postfix              // sorts the rule execution to the end
)
