package next

import (
	"testing"

	"github.com/ionous/iffy/rt"
)

func TestCompareNumbers(t *testing.T) {
	test := func(a float64, op CompareTo, b float64, res bool) {
		var run rt.Panic
		cmp := &CompareNum{&Number{a}, op, &Number{b}}

		if ok, e := cmp.GetBool(run); e != nil {
			t.Fatal(e)
		} else if res != ok {
			t.Fatal("mismatch")
		}
	}
	test(10, GreaterThan{}, 1, true)
	test(1, GreaterThan{}, 10, false)
	test(8, GreaterThan{}, 8, false)
	//
	test(10, LessThan{}, 1, false)
	test(1, LessThan{}, 10, true)
	test(8, LessThan{}, 8, false)
	//
	test(10, EqualTo{}, 1, false)
	test(1, EqualTo{}, 10, false)
	test(8, EqualTo{}, 8, true)
}
