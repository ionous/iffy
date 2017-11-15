package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestArgs(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testArgs(t, "", "", 0)) // arguments are optional.
	x = x && assert.NoError(testArgs(t, "a", "a", 1))
	x = x && assert.NoError(testArgs(t, "a b c", "a b c", 3))
	x = x && assert.NoError(testArgs(t, "a  b		c", "a b c", 3))
	x = x && assert.NoError(testArgs(t, "a b c  ", "a b c", 3))
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
	assert, x := testify.New(t), true
	x = x && assert.NoError(testArgx(t, "a {1+2}", "a 1 2 ADD"))
	x = x && assert.NoError(testArgx(t, "{(5+6)*(7+8)}", "5 6 ADD 7 8 ADD MUL"))
	x = x && assert.NoError(testArgx(t,
		"{{5|first!}+{'hello'|second! 6|third: 7}}",
		"5 FIRST/1 7 6 `hello` SECOND/2 THIRD/2 ADD",
	))
}

func testArgx(t *testing.T, str, want string) error {
	var p ArgParser
	return testRes(t, &p, str, want)
}
