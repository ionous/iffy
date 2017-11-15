package chart

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestPipe(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testPipe(t, "", "")) // arguments are optional.
	x = x && assert.NoError(testPipe(t, "world", "world"))
	x = x && assert.NoError(testPipe(t, "up! up", "up UP/1"))
	x = x && assert.NoError(testPipe(t, "up! up|bup! bup", "bup up UP/1 BUP/2"))
	x = x && assert.NoError(testPipe(t, "up!|up!", "UP/0 UP/1"))
	x = x && assert.NoError(testPipe(t, "up! up|bup!", "up UP/1 BUP/1"))
	x = x && assert.NoError(testPipe(t, "hello!", "HELLO/0"))
	x = x && assert.NoError(testPipe(t, "world|hello!", "world HELLO/1"))
	x = x && assert.NoError(testPipe(t, "world | hello! ", "world HELLO/1"))
	x = x && assert.NoError(testPipe(t, "world|hello! there", "there world HELLO/2"))
	x = x && assert.NoError(testPipe(t, "world|capitalize!|hello: there", "there world CAPITALIZE/1 HELLO/2"))
	x = x && assert.NoError(testPipe(t, "(5+6)*(7+8)", "5 6 ADD 7 8 ADD MUL"))
	//
	x = x && assert.Error(testPipe(t, "world|", ignoreResult))
	x = x && assert.Error(testPipe(t, "world|nofun", ignoreResult))
}

func testPipe(t *testing.T, str, want string) error {
	var p PipeParser
	return testRes(t, &p, str, want)
}
