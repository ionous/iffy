package chart

import (
	"testing"
)

func TestSeq(t *testing.T) {
	if e := testSeq(t, "", ""); e != nil {
		// arguments are optional.
		t.Fatal(e)
	}
	if e := testSeq(t, "a", "a"); e != nil {
		t.Fatal(e)
	}
	if e := testSeq(t, "x+y", "x y ADD"); e != nil {
		t.Fatal(e)
	}
	if e := testSeq(t, "x  +  y  ", "x y ADD"); e != nil {
		t.Fatal(e)
	}
	if e := testSeq(t, "(x+y)*z", "x y ADD z MUL"); e != nil {
		t.Fatal(e)
	}
	if e := testSeq(t, "( x + y ) * ( z ) ", "x y ADD z MUL"); e != nil {
		t.Fatal(e)
	}
	if e := testSeq(t, "(5+6)*(7+8)", "5 6 ADD 7 8 ADD MUL"); e != nil {
		t.Fatal(e)
	}
	if e := testSeq(t, "() ", ignoreResult); e == nil {
		t.Fatal(e)
	}
	if e := testSeq(t, "( x + y ) * () ", ignoreResult); e == nil {
		t.Fatal(e)
	}
}

func testSeq(t *testing.T, str, want string) error {
	var p SeriesParser
	return testRes(t, &p, str, want)
}
