package chart

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestSubdir(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testSub(t, "{fun!}", "FUN/0"))
	x = x && assert.NoError(testSub(t, "{args: a b}", "a b ARGS/2"))
	x = x && assert.NoError(testSub(t, "{quest?}", "QUEST/0"))
}

func testSub(t *testing.T, str, want string) error {
	var p SubdirParser
	return testRes(t, &p, str, want)
}
