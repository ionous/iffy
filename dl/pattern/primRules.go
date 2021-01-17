package pattern

import "github.com/ionous/iffy/rt"

// BoolRule responds with a true/false result when its filters are satisfied.
// It implements rt.BoolEval.
type xBoolRule struct {
	Filter rt.BoolEval
	rt.BoolEval
}

// NumberRule responds with a single number when its filters are satisfied.
type xNumberRule struct {
	Filter rt.BoolEval
	rt.NumberEval
}

// TextRule responds with a bit of text when its filters are satisfied.
type xTextRule struct {
	Filter rt.BoolEval
	rt.TextEval
}
