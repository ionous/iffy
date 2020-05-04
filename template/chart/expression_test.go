package chart

import (
	"testing"
)

func TestExpression(t *testing.T) {
	if e := testExp(t, "fun!", "FUN/0"); e != nil {
		t.Fatal(e)
	}
	if e := testExp(t, "call: a b", "a b CALL/2"); e != nil {
		t.Fatal(e)
	}
	if e := testExp(t, "quest?", "QUEST/0"); e != nil {
		t.Fatal(e)
	}
	if e := testExp(t, "x+y", "x y ADD"); e != nil {
		t.Fatal(e)
	}
	if e := testExp(t, "(5+6)*(7+8)", "5 6 ADD 7 8 ADD MUL"); e != nil {
		t.Fatal(e)
	}
	if e := testExp(t, "5*(6-4)", "5 6 4 SUB MUL"); e != nil {
		t.Fatal(e)
	}
	if e := testExp(t, "x and y", "x y LAND"); e != nil {
		t.Fatal(e)
	}
	if e := testExp(t, "a and (b or {isNot: c})", "a b c ISNOT/1 LOR LAND"); e != nil {
		t.Fatal(e)
	}
	if e := testExp(t, "!", ignoreResult); e == nil {
		t.Fatal(e)
	}
	if e := testExp(t, "fun!!", ignoreResult); e == nil {
		t.Fatal(e)
	}
}

func testExp(t *testing.T, str, want string) error {
	p := ExpressionParser{argFactory: &AnyFactory{}}
	return testRes(t, &p, str, want)
}

func testRes(t *testing.T, p ExpressionState, str, want string) (err error) {
	t.Logf("parsing: '%s'", str)
	if e := Parse(p, str); e != nil {
		t.Log("couldnt Parse", e)
		err = e
	} else if res, e := p.GetExpression(); e != nil {
		t.Log("invalid expression", e)
		err = e
	} else if want != ignoreResult {
		if got := res.String(); got != want {
			err = mismatched(want, got)
		} else {
			t.Log("ok", got)
		}
	}
	return
}

// for testing errors when we want to fail before the match is tested.
const ignoreResult = "~~IGNORE~~"
