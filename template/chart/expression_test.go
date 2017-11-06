package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestExpression(t *testing.T) {
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(testExp(t, "fun!", "FUN/0"))
	x = x && assert.NoError(testExp(t, "args: a b", "a b ARGS/2"))
	x = x && assert.NoError(testExp(t, "quest?", "QUEST/0"))
	x = x && assert.NoError(testExp(t, "x+y", "x y ADD"))
	x = x && assert.Error(testExp(t, "!", ignoreResult))
	x = x && assert.Error(testExp(t, "fun!!", ignoreResult))
}

func testExp(t *testing.T, str, want string) error {
	p := MakeExpressionParser(&AnyFactory{})
	return testRes(t, &p, str, want)
}

func testRes(t *testing.T, p ExpressionState, str, want string) (err error) {
	t.Logf("parsing: '%s'", str)
	if e := parse(p, str); e != nil {
		t.Log("couldnt parse", e)
		err = e
	} else if res, e := p.GetExpression(); e != nil {
		t.Log("invalid expression", e)
		err = e
	} else if want != ignoreResult {
		if got := res.String(); got != want {
			err = errutil.Fmt("want(%d): %s; != got(%d): %s.", len(want), want, len(got), got)
		} else {
			t.Log("ok", got)
		}
	}
	return
}

// for testing errors when we want to fail before the match is tested.
const ignoreResult = "~~IGNORE~~"
