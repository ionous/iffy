package pattern_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/test/testutil"
)

// TestFactorial of the number 3 to verify pattern recursion works.
func TestFactorial(t *testing.T) {
	// rules are run in reverse order.
	run := patternRuntime{PatternMap: testutil.PatternMap{
		"factorial": &pattern.Pattern{
			Name: "factorial",
			Params: []term.Preparer{
				&term.Number{Name: "num"},
			},
			Returns: &term.Number{Name: "num"},
			Rules: []*pattern.Rule{{
				Execute: core.NewActivity(
					&core.Assign{Var: N("num"),
						From: &core.FromNum{
							&core.ProductOf{
								A: V("num"),
								B: &pattern.Determine{
									Pattern: "factorial",
									Arguments: core.NamedArgs(
										"num", &core.FromNum{
											&core.DiffOf{
												V("num"),
												I(1),
											},
										},
									)}}}}),
			}, {
				Filter: &core.CompareNum{
					V("num"),
					&core.EqualTo{},
					I(0),
				},
				Execute: core.NewActivity(
					&core.Assign{Var: N("num"),
						From: &core.FromNum{
							I(1),
						}},
				),
			}}},
	}}
	// determine the factorial of the number 3
	det := pattern.Determine{
		Pattern: "factorial",
		Arguments: core.NamedArgs(
			"num", &core.FromNum{
				I(3),
			}),
	}
	if v, e := safe.GetNumber(&run, &det); e != nil {
		t.Fatal(e)
	} else if got, want := v.Int(), 3*(2*(1*1)); got != want {
		t.Fatal("mismatch: expected:", want, "have:", got)
	} else {
		t.Log("factorial okay", got)
	}
}
