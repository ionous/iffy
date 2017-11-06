package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestArg(t *testing.T) {
	test := func(str, want string, a int) (err error) {
		p := MakeArgParser(&AnyFactory{})
		t.Logf("parsing: '%s'", str)
		if e := parse(&p, str); e != nil {
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
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(test("", "", 0)) // arguments are optional.
	x = x && assert.NoError(test("a", "a", 1))
	x = x && assert.NoError(test("a b c", "a b c", 3))
	x = x && assert.NoError(test("a  b		c", "a b c", 3))
	x = x && assert.NoError(test("a b c  ", "a b c", 3))
}
