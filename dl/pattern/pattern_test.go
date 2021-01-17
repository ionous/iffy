package pattern_test

import (
	"fmt"

	"github.com/ionous/iffy/ephemera/debug"
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
		fmt.Printf(`say_me %d = "`, i)
		if e := debug.DetermineSay(i).Execute(&run); e != nil {
			fmt.Println("Error:", e)
		}
		fmt.Println(`"`)
	}

	// Output:
	// say_me 1 = "One!"
	// say_me 2 = "Two!"
	// say_me 3 = "Three!"
	// say_me 4 = "Not between 1 and 3."
}
