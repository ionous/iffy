package next

import "github.com/ionous/iffy/rt"

// Is transparently returns a boolean eval.
// It exists to help smooth the use of command expressions:
// eg. "is" {some expression}
type Is struct {
	rt.BoolEval
}

// IsNot returns the opposite of a boolean eval.
type IsNot struct {
	rt.BoolEval
}
