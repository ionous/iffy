package chart

import (
	"testing"
)

func TestPipe(t *testing.T) {
	if e := testPipe(t, "", ""); e != nil {
		// arguments are optional.
		t.Fatal(e)
	}
	if e := testPipe(t, ".world", "world"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, "up! .up", "up UP/1"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, "up! .up|bup! .bup", "bup up UP/1 BUP/2"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, "up!|up!", "UP/0 UP/1"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, "up! .up|bup!", "up UP/1 BUP/1"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, "hello!", "HELLO/0"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, ".world|hello!", "world HELLO/1"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, ".world | hello! ", "world HELLO/1"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, ".world|hello! .there", "there world HELLO/2"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, ".world|capitalize!|hello: .there", "there world CAPITALIZE/1 HELLO/2"); e != nil {
		t.Fatal(e)
	}
	if e := testPipe(t, "(5+6)*(7+8)", "5 6 ADD 7 8 ADD MUL"); e != nil {
		t.Fatal(e)
	}
	//
	if e := testPipe(t, ".world|", ignoreResult); e == nil {
		t.Fatal(e)
	} else {
		t.Log("ok, error:", e)
	}
	if e := testPipe(t, ".world|.nofun", ignoreResult); e == nil {
		t.Fatal(e)
	} else {
		t.Log("ok, error:", e)
	}
}

func testPipe(t *testing.T, str, want string) error {
	var p PipeParser
	return testRes(t, &p, str, want)
}
