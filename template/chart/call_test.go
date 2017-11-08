package chart

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestCall(t *testing.T) {
	assert, x := testify.New(t), true
	f := &AnyFactory{}
	x = x && assert.NoError(testCall(t, f, "fun!", "FUN/0"))
	x = x && assert.NoError(testCall(t, f, "call: a b", "a b CALL/2"))
	x = x && assert.NoError(testCall(t, f, "quest?", "QUEST/0"))
}

func TestCallSubdir(t *testing.T) {
	assert, x := testify.New(t), true
	var f ExpressionStateFactory
	x = x && assert.NoError(testCall(t, f, "call: {5+6}", "5 6 ADD CALL/1"))
	x = x && assert.NoError(testCall(t, f, "call: {a!} b", "A/0 b CALL/2"))
	x = x && assert.NoError(testCall(t, f, "call: a {1+2}", "a 1 2 ADD CALL/2"))
}

func TestCallSubSubdir(t *testing.T) {
	assert, x := testify.New(t), true
	var f ExpressionStateFactory
	x = x && assert.NoError(testCall(t, f,
		"call: {{5|first!}+{'hello'|second! 6|third: 7}}",
		"5 FIRST/1 7 6 `hello` SECOND/2 THIRD/2 ADD CALL/1",
	))
}

func testCall(t *testing.T, f ExpressionStateFactory, str string, want string) error {
	p := MakeCallParser(0, f)
	return testRes(t, &p, str, want)
}
