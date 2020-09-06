package pattern_test

import (
	"fmt"

	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// ExampleSayMe converts numbers to text
// http://learnyouahaskell.com/syntax-in-functions
func ExampleSayMe() {
	// rules are run in reverse order.
	run := patternRuntime{patternMap: patternMap{
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
	scope.ScopeStack // parameters are pushed onto the stack.
	patternMap       // holds pointers to patterns
}
type patternMap map[string]interface{}

// skip assembling the pattern from the db
// we just want to test we can invoke a pattern successfully.
// pv is a pointer to a pattern instance, and we copy its contents in.
func (m *patternRuntime) GetEvalByName(name string, pv interface{}) (err error) {
	if patternPtr, ok := m.patternMap[name]; ok {
		stored := r.ValueOf(patternPtr).Elem()
		outVal := r.ValueOf(pv).Elem()
		outVal.Set(stored)

	} else {
		err = errutil.New("patternRuntime: unknown pattern", name)
	}
	return
}
