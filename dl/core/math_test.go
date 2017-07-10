package core_test

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
)

func (assert *CoreSuite) TestAdd() {
	assert.matchLine("11", func(c *ops.Builder) {
		if c.Cmd("print num").Begin() {
			c.Cmd("add", 1, 10)
			c.End()
		}
	})
}

func (assert *CoreSuite) TestSubtract() {
	assert.matchLine("-9", func(c *ops.Builder) {
		if c.Cmd("print num").Begin() {
			c.Cmd("sub", 1, 10)
			c.End()
		}
	})
}

func (assert *CoreSuite) TestMultiply() {
	assert.matchLine("200", func(c *ops.Builder) {
		if c.Cmd("print num").Begin() {
			c.Cmd("mul", 20, 10)
			c.End()
		}
	})
}

// TestDivide tests numbers directly.
func (assert *CoreSuite) TestDivide() {
	var root struct{ Eval rt.NumberEval }
	if c, ok := assert.ops.NewBuilder(&root); ok {
		c.Cmd("div", 10, 2)
		//
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if v, e := root.Eval.GetNumber(run); assert.NoError(e) {
				assert.EqualValues(5, v)
			}
		}
	}
}

// TestDivideByZero should not panic, but simply error.
func (assert *CoreSuite) TestDivideByZero() {
	var root struct{ Eval rt.NumberEval }
	if c, ok := assert.ops.NewBuilder(&root); ok {
		c.Cmd("div", 10, 0)
		//
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if _, e := root.Eval.GetNumber(run); assert.Error(e) {
			}
		}
	}

}
