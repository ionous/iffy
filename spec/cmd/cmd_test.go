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
	if c, ok := cmd.NewBuilder(); assert.True(ok) {
		if c.Cmds().Block() {
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

// TestDuplicateKeys verifies duplicate keys results in error.
func TestDuplicateKeys(t *testing.T) {
	assert := assert.New(t)
	if c, ok := cmd.NewBuilder(); assert.True(ok) {
		c.Param("unique").Cmd("unit")
		c.Param("unique").Cmd("unit")
		if root, e := c.Build(); !assert.Error(e) {
			cmd.Print(root)
		}
	}
}

// TestSpec verifies that a valid spec results in no error.
func TestSpec(t *testing.T) {
	assert := assert.New(t)
	if c, ok := cmd.NewBuilder(); assert.True(ok) {
		if c.Cmd("unit").Block() {
			if c.Param("trials").Cmds().Block() {
				// 	// cycles:
				c.Cmd("match output", sliceOf.String("a", "b", "c", "d"))
				if c.Cmd("match output", sliceOf.String("a", "b", "c", "d")).Block() {
					if c.Cmd("for each num", sliceOf.Float(1, 2, 3, 4)).Block() {
						a := c.Cmd("cycle", sliceOf.String("a", "b", "c"))
						c.Cmd("print text", a).End()
					}
					c.End()
				}
				// stopping:
				if c.Cmd("match output", sliceOf.String("a", "b", "c", "c")).Block() {
					if c.Cmd("for each num", sliceOf.Float(1, 2, 3, 4)).Block() {
						c.Cmd("print text", c.Cmd("stopping", sliceOf.String("a", "b", "c")))
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
