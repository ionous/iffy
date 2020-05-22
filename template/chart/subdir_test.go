package chart

import (
	"testing"
)

func TestSubdir(t *testing.T) {
	if e := testSub(t, "{fun!}", "FUN/0"); e != nil {
		t.Fatal(e)
	}
	if e := testSub(t, "{call: .a .b}", "a b CALL/2"); e != nil {
		t.Fatal(e)
	}
	if e := testSub(t, "{quest?}", "QUEST/0"); e != nil {
		t.Fatal(e)
	}
	if e := testSub(t, "{(5+6)*(7+8)}", "5 6 ADD 7 8 ADD MUL"); e != nil {
		t.Fatal(e)
	}
}

func testSub(t *testing.T, str, want string) error {
	var p SubdirParser
	return testRes(t, &p, str, want)
}
