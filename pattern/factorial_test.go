package pattern_test

import (
	"testing"

	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/dl/next"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

// TestFactorial of the number 3 to verify pattern recursion works.
func TestFactorial(t *testing.T) {
	// rules are run in reverse order.
	run := patternRuntime{patternMap: patternMap{
		"factorial": pattern.NumberRules{
			{
				NumberEval: &next.ProductOf{
					&next.GetVar{"num"},
					&next.Determine{
						"factorial",
						scope.Parameters{
							"num": &next.DiffOf{
								&next.GetVar{"num"},
								&next.Number{1},
							},
						},
					},
				},
			}, {
				Filters: []rt.BoolEval{
					&next.CompareNum{
						&next.GetVar{"num"},
						&next.EqualTo{},
						&next.Number{0},
					},
				},
				NumberEval: &next.Number{1},
			},
		}}}
	// determine the factorial of the number 3
	det := next.Determine{"factorial", scope.Parameters{
		"num": &next.Number{3},
	}}
	if v, e := rt.GetNumber(&run, &det); e != nil {
		t.Fatal(e)
	} else if want := 3.0 * (2 * (1 * 1)); v != want {
		t.Fatal(assign.Mismatch(t.Name(), want, v))
	} else {
		t.Log("factorial okay", v)
	}
}
