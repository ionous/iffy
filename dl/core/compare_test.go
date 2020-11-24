package core

import (
	"testing"

	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/test/testutil"
)

func TestCompareNumbers(t *testing.T) {
	test := func(a float64, op Comparator, b float64, res bool) {
		var run testutil.PanicRuntime
		cmp := &CompareNum{&Number{a}, op, &Number{b}}
		if ok, e := safe.GetBool(run, cmp); e != nil {
			t.Fatal(e)
		} else if res != ok.Bool() {
			t.Fatal("mismatch")
		}
	}
	test(10, &GreaterThan{}, 1, true)
	test(1, &GreaterThan{}, 10, false)
	test(8, &GreaterThan{}, 8, false)
	//
	test(10, &LessThan{}, 1, false)
	test(1, &LessThan{}, 10, true)
	test(8, &LessThan{}, 8, false)
	//
	test(10, &EqualTo{}, 1, false)
	test(1, &EqualTo{}, 10, false)
	test(8, &EqualTo{}, 8, true)
}

func TestCompareText(t *testing.T) {
	test := func(a string, op Comparator, b string, res bool) {
		var run testutil.PanicRuntime
		cmp := &CompareText{&Text{a}, op, &Text{b}}
		if ok, e := safe.GetBool(run, cmp); e != nil {
			t.Fatal(e)
		} else if res != ok.Bool() {
			t.Fatal("mismatch")
		}
	}
	test("Z", &GreaterThan{}, "A", true)
	test("A", &GreaterThan{}, "Z", false)
	//
	test("marzip", &LessThan{}, "marzipan", true)
	test("marzipan", &LessThan{}, "marzip", false)
	//
	test("bobby", &EqualTo{}, "bobby", true)
	test("bobby", &EqualTo{}, "phillipa", false)
}
