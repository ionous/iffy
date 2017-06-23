package cmd_test

import (
	"github.com/ionous/iffy/spec/cmd"
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestMismatchedEnds verifies that a mismatched end results in an error
func TestMismatchedEnds(t *testing.T) {
	assert := assert.New(t)
	if c, ok := cmd.NewBuilder(); ok {
		if c.Array().Block() {
			if c.Cmd("unit").Block() {
				// c.End() --> missing an end
			} // unit
			c.End()
		} // array block
		if root, e := c.Build(); !assert.Error(e) {
			cmd.Print(root)
		}
	}
}

// TestSpec verifies that a valid spec results in no error.
func TestSpec(t *testing.T) {
	assert := assert.New(t)
	if c, ok := cmd.NewBuilder(); ok {
		if c.Cmd("unit").Block() {
			if c.Param("trials").Array().Block() {
				// cycles
				if c.Cmd("match output", c.Val(sliceOf.String("a", "b", "c", "d"))).Block() {
					if c.Cmd("for each num", c.Val(sliceOf.Float(1, 2, 3, 4))).Block() {
						c.Cmd("print text", c.Cmd("cycle", c.Val(sliceOf.String("a", "b", "c"))))
						c.End()
					}
					c.End()
				}
				// stopping
				if c.Cmd("match output", c.Val(sliceOf.String("a", "b", "c", "c"))).Block() {
					if c.Cmd("for each num", c.Val(sliceOf.Float(1, 2, 3, 4))).Block() {
						c.Cmd("print text", c.Cmd("stopping", c.Val(sliceOf.String("a", "b", "c"))))
						c.End()
					}
					c.End()
				}
				c.End()
			} // trials
			c.End()
		} // unit
		if root, e := c.Build(); assert.NoError(e) {
			cmd.Print(root)
		}
	}
}
