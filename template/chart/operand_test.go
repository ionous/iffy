package chart

import (
	"testing"

	"github.com/ionous/iffy/template/types"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
)

func TestOperand(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testOp(t, "'hello'", types.Quote("hello").String()))
	x = x && assert.NoError(testOp(t, "1.2", types.Number(1.2).String()))
	x = x && assert.NoError(testOp(t, "true", types.Bool(true).String()))
	x = x && assert.NoError(testOp(t, "false", types.Bool(false).String()))
	x = x && assert.NoError(testOp(t, "object", types.Reference(sliceOf.String("object")).String()))
	x = x && assert.NoError(testOp(t, "a", types.Reference(sliceOf.String("a")).String()))
	x = x && assert.NoError(testOp(t, "object.property", types.Reference(sliceOf.String("object", "property")).String()))
	x = x && assert.Error(testOp(t, "#", ignoreResult))
}

func testOp(t *testing.T, str, want string) error {
	var p OperandParser
	return testRes(t, &p, str, want)
}
