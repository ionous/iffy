package chart

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestCall(t *testing.T) {
	assert, x := testify.New(t), true
	f := &AnyFactory{}
	x = x && assert.NoError(testCall(t, f, "fun!", "FUN/0"))
	x = x && assert.NoError(testCall(t, f, "args: a b", "a b ARGS/2"))
	x = x && assert.NoError(testCall(t, f, "quest?", "QUEST/0"))
}

func TestCallDir(t *testing.T) {
	assert, x := testify.New(t), true
	var f ExpressionStateFactory
	x = x && assert.NoError(testCall(t, f, "args: {5+6}", "5 6 ADD ARGS/1"))
	x = x && assert.NoError(testCall(t, f, "args: {a!} b", "A/0 b ARGS/2"))
	x = x && assert.NoError(testCall(t, f, "args: a {1+2}", "a 1 2 ADD ARGS/2"))
}

func testCall(t *testing.T, f ExpressionStateFactory, str string, want string) error {
	p := MakeCallParser(0, f)
	return testRes(t, &p, str, want)
}
