package pattern_test

import (
	"fmt"

	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/test/testutil"
)

// ExampleSayMe converts numbers to text
// http://learnyouahaskell.com/syntax-in-functions
func ExampleSayMe() {
	// rules are run in reverse order.
	run := patternRuntime{PatternMap: testutil.PatternMap{
		"say_me": &debug.SayPattern,
	}}
	// say 4 numbers
	for i := 1; i <= 4; i++ {
		if text, e := debug.DetermineSay(i).GetText(&run); e != nil {
			fmt.Println("Error:", e)
			break
		} else {
			fmt.Println(fmt.Sprintf("say_me %d = \"%s\"", i, text))
		}
	}

	// Output:
	// say_me 1 = "One!"
	// say_me 2 = "Two!"
	// say_me 3 = "Three!"
	// say_me 4 = "Not between 1 and 3."
}

type baseRuntime struct {
	testutil.PanicRuntime
}

type patternRuntime struct {
	baseRuntime
	scope.ScopeStack    // parameters are pushed onto the stack.
	testutil.PatternMap // holds pointers to patterns
}
