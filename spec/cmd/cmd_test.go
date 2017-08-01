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
		if c.Cmds().Begin() {
			if c.Cmd("unit").Begin() {
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
		if c.Cmd("unit").Begin() {
			if c.Param("trials").Cmds().Begin() {
				// 	// cycles:
				c.Cmd("match output", sliceOf.String("a", "b", "c", "d"))
				if c.Cmd("match output", sliceOf.String("a", "b", "c", "d")).Begin() {
					if c.Cmd("for each", sliceOf.Float(1, 2, 3, 4)).Begin() {
						a := c.Cmd("cycle", sliceOf.String("a", "b", "c"))
						c.Cmd("print text", a).End()
					}
					c.End()
				}
				// stopping:
				if c.Cmd("match output", sliceOf.String("a", "b", "c", "c")).Begin() {
					if c.Cmd("for each", sliceOf.Float(1, 2, 3, 4)).Begin() {
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

var order = &cmd.Command{
	Name: "container",
	Args: []interface{}{
		5,
		&cmd.Command{
			Name: "op",
		},
		10,
	},
}

func TestPositioning(t *testing.T) {
	assert := assert.New(t)
	if c, ok := cmd.NewBuilder(); assert.True(ok) {
		c.Cmd("container", c.Val(5), c.Cmd("op"), c.Val(10))
		if root, e := c.Build(); assert.NoError(e) {
			assert.EqualValues(order, root.Args[0])
		}
	}
}

func TestChaining(t *testing.T) {
	assert := assert.New(t)
	if c, ok := cmd.NewBuilder(); assert.True(ok) {

		if c.Cmd("container").Begin() {
			c.Val(5).Cmd("op").Val(10)
			c.End()
		}
		if root, e := c.Build(); assert.NoError(e) {
			assert.EqualValues(order, root.Args[0])
		}
	}
}

func TestParameterChaining(t *testing.T) {
	assert := assert.New(t)
	if c, ok := cmd.NewBuilder(); assert.True(ok) {
		assert.Panics(func() {
			c.Cmd("container", c.Val(5).Val(10))
		})
	}
}
