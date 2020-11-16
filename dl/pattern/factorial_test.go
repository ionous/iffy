package pattern_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
)

// TestFactorial of the number 3 to verify pattern recursion works.
func TestFactorial(t *testing.T) {
	// rules are run in reverse order.
	run := patternRuntime{PatternMap: pattern.PatternMap{
		"factorial": &pattern.NumberPattern{
			pattern.CommonPattern{
				Name: "factorial",
				Prologue: []term.Preparer{
					&term.Number{Name: "num"},
				},
			}, []*pattern.NumberRule{{
				NumberEval: &core.ProductOf{
					&core.GetVar{Name: &core.Text{"num"}},
					&pattern.DetermineNum{
						Pattern: "factorial", Arguments: pattern.NewNamedParams(
							"num", &core.FromNum{
								&core.DiffOf{
									&core.GetVar{Name: &core.Text{"num"}},
									&core.Number{1},
								},
							},
						)}},
			}, {
				Filter: &core.CompareNum{
					&core.GetVar{Name: &core.Text{"num"}},
					&core.EqualTo{},
					&core.Number{0},
				},
				NumberEval: &core.Number{1},
			}}},
	}}
	// determine the factorial of the number 3
	det := pattern.DetermineNum{
		Pattern: "factorial", Arguments: pattern.NewNamedParams(
			"num", &core.FromNum{
				&core.Number{3},
			}),
	}
	if v, e := rt.GetNumber(&run, &det); e != nil {
		t.Fatal(e)
	} else if want := 3.0 * (2 * (1 * 1)); v != want {
		t.Fatal("mismatch: expected:", want, "have:", v)
	} else {
		t.Log("factorial okay", v)
	}
}
