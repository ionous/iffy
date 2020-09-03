package pattern

import "github.com/ionous/iffy/rt"

// BoolRule responds with a true/false result when its filters are satisfied.
// It implements rt.BoolEval.
type BoolRule struct {
	Filter rt.BoolEval
	rt.BoolEval
}

// NumberRule responds with a single number when its filters are satisfied.
type NumberRule struct {
	Filter rt.BoolEval
	rt.NumberEval
}

// TextRule responds with a bit of text when its filters are satisfied.
type TextRule struct {
	Filter rt.BoolEval
	rt.TextEval
}
