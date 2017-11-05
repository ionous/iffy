package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestArg(t *testing.T) {
	test := func(str string, want string) (err error) {
		p := MakeArgParser(&AnyFactory{})
		t.Logf("parsing: '%s'", str)
		if e := parse(&p, str); e != nil {
			err = e
		} else if res, arity, e := p.GetArguments(); e != nil {
			err = e
		} else if a := strings.Count(want, "["); arity != a {
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
	x = x && assert.NoError(test("", "")) // arguments are optional.
	x = x && assert.NoError(test("a", "[a]"))
	x = x && assert.NoError(test("a b c", "[a][b][c]"))
	x = x && assert.NoError(test("a  b		c", "[a][b][c]"))
	x = x && assert.NoError(test("a b c  ", "[a][b][c]"))
}
