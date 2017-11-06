package chart

import (
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestOperand(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testOp(t, "'hello'", Quote("hello").String()))
	x = x && assert.NoError(testOp(t, "1.2", Number(1.2).String()))
	x = x && assert.NoError(testOp(t, "object", Reference(sliceOf.String("object")).String()))
	x = x && assert.NoError(testOp(t, "a", Reference(sliceOf.String("a")).String()))
	x = x && assert.NoError(testOp(t, "object.property", Reference(sliceOf.String("object", "property")).String()))
	x = x && assert.Error(testOp(t, "#", ignoreResult))
}

func testOp(t *testing.T, str, want string) error {
	var p OperandParser
	return testRes(t, &p, str, want)
}
