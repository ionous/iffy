package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestExp(t *testing.T) {
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(testExp(t, "fun!", "FUN/0"))
	x = x && assert.NoError(testExp(t, "args: a b", "[a][b]ARGS/2"))
	x = x && assert.NoError(testExp(t, "quest?", "QUEST/0"))
	x = x && assert.NoError(testExp(t, "x+y", "[x][y]ADD"))
	x = x && assert.Error(testExp(t, "!", ignoreResult))
	x = x && assert.Error(testExp(t, "fun!!", ignoreResult))
}

func testExp(t *testing.T, str string, want string) (err error) {
	p := MakeExpressionParser(&AnyFactory{})
	t.Logf("parsing: '%s'", str)
	if e := parse(&p, str); e != nil {
		t.Log("couldnt parse", e)
		err = e
	} else if res, e := p.GetExpression(); e != nil {
		t.Log("invalid expression", e)
		err = e
	} else if want != ignoreResult {
		got := res.String()
		t.Log("got", got)
		if got != want {
			err = errutil.New("want", want)
		}
	}
	return
}

// for testing errors when we want to fail before the match is tested.
const ignoreResult = "~~IGNORE~~"
