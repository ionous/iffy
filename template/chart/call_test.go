package chart

import (
	"testing"
)

func TestCall(t *testing.T) {
	f := &AnyFactory{}

	if e := testCall(t, f, "call: .a .b", "a b CALL/2"); e != nil {
		t.Fatal(e)
	}
	if e := testCall(t, f, "fun!", "FUN/0"); e != nil {
		t.Fatal(e)
	}
	if e := testCall(t, f, "quest?", "QUEST/0"); e != nil {
		t.Fatal(e)
	}
}

func TestCallSubdir(t *testing.T) {
	var f ExpressionStateFactory
	if e := testCall(t, f, "call: {5+6}", "5 6 ADD CALL/1"); e != nil {
		t.Fatal(e)
	}
	if e := testCall(t, f, "call: {a!} .b", "A/0 b CALL/2"); e != nil {
		t.Fatal(e)
	}
	if e := testCall(t, f, "call: .a {1+2}", "a 1 2 ADD CALL/2"); e != nil {
		t.Fatal(e)
	}
}

func TestCallSubSubdir(t *testing.T) {
	var f ExpressionStateFactory
	if e := testCall(t, f,
		"call: {{5|first!}+{'hello'|second! 6|third: 7}}",
		"5 FIRST/1 7 6 \"hello\" SECOND/2 THIRD/2 ADD CALL/1",
	); e != nil {
		t.Fatal(e)
	}
}

func testCall(t *testing.T, f ExpressionStateFactory, str string, want string) error {
	p := MakeCallParser(0, f)
	return testRes(t, &p, str, want)
}
