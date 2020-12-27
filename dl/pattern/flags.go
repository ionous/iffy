package pattern

// Flags controlling patterns capable of producing multiple results.
// Rules satisfying their filters get sorted and processed in the flagged order.
// see also: splitRules
type Flags int

//go:generate stringer -type=Flags
const (
	Terminal Flags = iota // default, stops considering other rules
	Prefix                // sorts the rule execution to the front
	Postfix               // sorts the rule execution to the end
)
