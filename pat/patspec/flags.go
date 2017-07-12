package patspec

import (
	"github.com/ionous/iffy/pat"
)

type PatternTiming interface {
	Flags() pat.Flags
}

// ContinueAfter tells the patten matcher to keep going after running the current pattern. this...others.
type ContinueAfter struct{}

// ContinueBefore tells the pattern matcher to run other matching patterns, and then run this pattern. other...this.
type ContinueBefore struct{}

func (ContinueAfter) Flags() pat.Flags  { return pat.Prefix }
func (ContinueBefore) Flags() pat.Flags { return pat.Postfix }
