package chart

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestSubdir(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testSub(t, "{fun!}", "FUN/0"))
	x = x && assert.NoError(testSub(t, "{call: a b}", "a b CALL/2"))
	x = x && assert.NoError(testSub(t, "{quest?}", "QUEST/0"))
	x = x && assert.NoError(testSub(t, "{(5+6)*(7+8)}", "5 6 ADD 7 8 ADD MUL"))
}

func testSub(t *testing.T, str, want string) error {
	var p SubdirParser
	return testRes(t, &p, str, want)
}
