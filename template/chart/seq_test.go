package chart

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestSeq(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testSeq(t, "", "")) // arguments are optional.
	x = x && assert.NoError(testSeq(t, "a", "a"))
	x = x && assert.NoError(testSeq(t, "x+y", "x y ADD"))
	x = x && assert.NoError(testSeq(t, "x  +  y  ", "x y ADD"))
	x = x && assert.NoError(testSeq(t, "(x+y)*z", "x y ADD z MUL"))
	x = x && assert.NoError(testSeq(t, "( x + y ) * ( z ) ", "x y ADD z MUL"))
	x = x && assert.NoError(testSeq(t, "(5+6)*(7+8)", "5 6 ADD 7 8 ADD MUL"))
	x = x && assert.Error(testSeq(t, "() ", ignoreResult))
	x = x && assert.Error(testSeq(t, "( x + y ) * () ", ignoreResult))
}

func testSeq(t *testing.T, str, want string) error {
	var p SequenceParser
	return testRes(t, &p, str, want)
}
