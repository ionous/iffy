package pattern_test

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/dl/next"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

func TestFactorial(t *testing.T) {
	var run patternRuntime
	if v, e := rt.GetNumber(&run, &next.Determine{"factorial", scope.Parameters{
		"num": &next.Number{3},
	}}); e != nil {
		t.Fatal(e)
	} else if want := 3.0 * (2 * (1 * 1)); v != want {
		t.Fatal(assign.Mismatch(t.Name(), want, v))
	} else {
		t.Log("factorial okay", v)
	}

}

type baseRuntime struct {
	rt.Panic
}

type patternRuntime struct {
	baseRuntime
	scope.ScopeStack // parameters are pushed onto the stack.
}

// skip assembling the pattern from the db
// we just want to test we can invoke a pattern successfully.
func (m *patternRuntime) GetField(name, field string) (ret interface{}, err error) {
	switch field {
	case object.Pattern:
		switch name {
		case "factorial":
			ret = factorial
		default:
			err = errutil.New("unknown pattern", field)
		}
	default:
		err = errutil.New("unknown field", field)
	}
	return
}

var factorial rt.NumberEval = pattern.NumberRules{
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
}
