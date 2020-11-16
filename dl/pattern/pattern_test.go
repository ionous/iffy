package pattern_test

import (
	"fmt"

	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// ExampleSayMe converts numbers to text
// http://learnyouahaskell.com/syntax-in-functions
func ExampleSayMe() {
	// rules are run in reverse order.
	run := patternRuntime{PatternMap: pattern.PatternMap{
		"sayMe": &debug.SayPattern,
	}}
	// say 4 numbers
	for i := 1; i <= 4; i++ {
		if text, e := debug.DetermineSay(i).GetText(&run); e != nil {
			fmt.Println("Error:", e)
			break
		} else {
			fmt.Println(fmt.Sprintf("sayMe %d = \"%s\"", i, text))
		}
	}

	// Output:
	// sayMe 1 = "One!"
	// sayMe 2 = "Two!"
	// sayMe 3 = "Three!"
	// sayMe 4 = "Not between 1 and 3."
}

type baseRuntime struct {
	rt.Panic
}

type patternRuntime struct {
	baseRuntime
	scope.ScopeStack   // parameters are pushed onto the stack.
	pattern.PatternMap // holds pointers to patterns
}
