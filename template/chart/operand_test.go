package chart

import (
	"testing"

	"github.com/ionous/iffy/template/types"
	"github.com/ionous/sliceOf"
)

func TestOperand(t *testing.T) {
	if e := testOp(t, "'hello'", types.Quote("hello").String()); e != nil {
		t.Fatal(e)
	}
	if e := testOp(t, "1.2", types.Number(1.2).String()); e != nil {
		t.Fatal(e)
	}
	if e := testOp(t, "true", types.Bool(true).String()); e != nil {
		t.Fatal(e)
	}
	if e := testOp(t, "false", types.Bool(false).String()); e != nil {
		t.Fatal(e)
	}
	if e := testOp(t, "object", types.Reference(sliceOf.String("object")).String()); e != nil {
		t.Fatal(e)
	}
	if e := testOp(t, "a", types.Reference(sliceOf.String("a")).String()); e != nil {
		t.Fatal(e)
	}
	if e := testOp(t, "object.property", types.Reference(sliceOf.String("object", "property")).String()); e != nil {
		t.Fatal(e)
	}
	if e := testOp(t, "#", ignoreResult); e == nil {
		t.Fatal(e)
	}
}

func testOp(t *testing.T, str, want string) error {
	var p OperandParser
	return testRes(t, &p, str, want)
}
