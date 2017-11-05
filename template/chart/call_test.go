package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestCall(t *testing.T) {
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(testCall(t, "fun!", "FUN/0"))
	x = x && assert.NoError(testCall(t, "args: a b", "[a][b]ARGS/2"))
	x = x && assert.NoError(testCall(t, "quest?", "QUEST/0"))
}

func testCall(t *testing.T, str string, want string) (err error) {
	p := MakeCallParser(0, &AnyFactory{})
	t.Logf("parsing: '%s'", str)
	if e := parse(&p, str); e != nil {
		err = e
	} else if res, e := p.GetExpression(); e != nil {
		err = e
	} else {
		got := res.String()
		t.Log("got:", got)
		if got != want {
			err = errutil.New("want", want)
		}
	}
	return
}
