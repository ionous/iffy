package chart

import (
	"github.com/ionous/iffy/template"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestOperand(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testOp(t, "'hello'", template.Quote("hello").String()))
	x = x && assert.NoError(testOp(t, "1.2", template.Number(1.2).String()))
	x = x && assert.NoError(testOp(t, "object", template.Reference(sliceOf.String("object")).String()))
	x = x && assert.NoError(testOp(t, "a", template.Reference(sliceOf.String("a")).String()))
	x = x && assert.NoError(testOp(t, "object.property", template.Reference(sliceOf.String("object", "property")).String()))
	x = x && assert.Error(testOp(t, "#", ignoreResult))
}

func testOp(t *testing.T, str, want string) error {
	var p OperandParser
	return testRes(t, &p, str, want)
}
