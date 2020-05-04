package pattern_test

import (
	"fmt"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// ExampleSayMe converts numbers to text
// http://learnyouahaskell.com/syntax-in-functions
func ExampleSayMe() {
	// rules are run in reverse order.
	run := patternRuntime{patternMap: patternMap{
		"sayMe": pattern.TextRules{
			{nil, SayIt("Not between 1 and 3")},
			{matchNumber(3), SayIt("San!")},
			{matchNumber(3), SayIt("Three!")},
			{matchNumber(2), SayIt("Two!")},
			{matchNumber(1), SayIt("One!")},
		}}}

	// say 4 numbers
	for i := 1; i <= 4; i++ {
		det := core.DetermineText{
			"sayMe", &core.Parameters{[]*core.Parameter{{
				"num",
				&core.FromNum{
					&core.Number{float64(i)},
				},
			}}},
		}
		if text, e := rt.GetText(&run, &det); e != nil {
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
	// sayMe 4 = "Not between 1 and 3"
}

type baseRuntime struct {
	rt.Panic
}

type patternRuntime struct {
	baseRuntime
	scope.ScopeStack // parameters are pushed onto the stack.
	patternMap
}
type patternMap map[string]interface{}

// skip assembling the pattern from the db
// we just want to test we can invoke a pattern successfully.
func (m *patternRuntime) GetField(name, field string) (ret interface{}, err error) {
	switch field {
	case object.Pattern:
		if p, ok := m.patternMap[name]; ok {
			ret = p
		} else {
			err = errutil.New("unknown pattern", field)
		}
	default:
		err = errutil.New("unknown field", field)
	}
	return
}

func SayIt(s string) rt.TextEval {
	return &core.Text{s}
}

type matchNumber int

func (m matchNumber) GetBool(run rt.Runtime) (okay bool, err error) {
	if v, e := run.GetVariable("num"); e != nil {
		err = e
	} else {
		n := int(v.(float64))
		okay = n == int(m)
	}
	return
}
