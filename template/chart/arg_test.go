package chart

import (
	"testing"

	"github.com/ionous/errutil"
)

func TestArgs(t *testing.T) {
	if e := testArgs(t, "", "", 0); e != nil {
		t.Fatal(e, "arguments are optional")
	}
	if e := testArgs(t, "a", "a", 1); e != nil {
		t.Fatal(e)
	}
	if e := testArgs(t, "a b c", "a b c", 3); e != nil {
		t.Fatal(e)
	}
	if e := testArgs(t, "a  b		c", "a b c", 3); e != nil {
		t.Fatal(e)
	}
	if e := testArgs(t, "a b c  ", "a b c", 3); e != nil {
		t.Fatal(e)
	}
}

func testArgs(t *testing.T, str, want string, a int) (err error) {
	p := MakeArgParser(&AnyFactory{})
	t.Logf("parsing: '%s'", str)
	if e := Parse(&p, str); e != nil {
		err = e
	} else if res, arity, e := p.GetArguments(); e != nil {
		err = e
	} else if a != arity {
		err = errutil.New("mismatched arity", arity, a)
	} else {
		got := res.String()
		t.Log("got:", got)
		if got != want {
			err = errutil.New("wanted:", want)
		}
	}
	return
}

func TestArgExpression(t *testing.T) {
	if e := testArgx(t, "a {1+2}", "a 1 2 ADD"); e != nil {
		t.Fatal(e)
	}
	if e := testArgx(t, "{(5+6)*(7+8)}", "5 6 ADD 7 8 ADD MUL"); e != nil {
		t.Fatal(e)
	}
	if e := testArgx(t,
		"{{5|first!}+{'hello'|second! 6|third: 7}}",
		"5 FIRST/1 7 6 \"hello\" SECOND/2 THIRD/2 ADD",
	); e != nil {
		t.Fatal(e)
	}
}

func testArgx(t *testing.T, str, want string) error {
	var p ArgParser
	return testRes(t, &p, str, want)
}
