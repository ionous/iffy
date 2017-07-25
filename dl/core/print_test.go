package core_test

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/sliceOf"
)

// TestPrintSpacing verifies that print span adds spaces between words as needed.
func (assert *CoreSuite) TestPrintSpacing() {
	var root struct {
		Eval rt.Execute
	}
	if c, ok := assert.ops.NewBuilder(&root); ok {
		if c.Cmd("print span").Begin() {
			if c.Cmds().Begin() {
				c.Cmd("print text", "hello")
				c.Cmd("print text", "there,")
				c.Cmd("print text", "world.")
				c.End()
			}
			c.End()
		}
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if e := root.Eval.Execute(run); assert.NoError(e) {
				lines := assert.Lines()
				assert.Equal(sliceOf.String("hello there, world."), lines, "Note the space after the comma, and the lack of space after the period.")
			}
		}
	}
}

func (assert *CoreSuite) TestPrintNum() {
	var root struct {
		Eval rt.Execute
	}
	if c, ok := assert.ops.NewBuilder(&root); ok {
		if c.Cmd("print span").Begin() {
			if c.Cmds().Begin() {
				c.Cmd("print num", 213)
				c.Cmd("print num word", 213)
				c.End()
			}
			c.End()
		}
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if e := root.Eval.Execute(run); assert.NoError(e) {
				lines := assert.Lines()
				assert.Equal(sliceOf.String("213 two hundred thirteen"), lines)
			}
		}
	}
}

// TestMultiLines verifies that iffy printing works similar to sashimi printing.
// In sashimi, the default printer made every print a new line, we should do the same. This test complements TestSingleLines.
func (assert *CoreSuite) TestMultiLines() {
	var root struct {
		Eval rt.Execute
	}
	if c, ok := assert.ops.NewBuilder(&root); ok {
		if c.Cmd("for each text").Begin() {
			c.Param("in").Val(sliceOf.String("hello", "there", "world"))
			if c.Param("go").Cmds().Begin() {
				c.Cmd("print text", c.Cmd("get", "@", "text"))
				c.End()
			}
			c.End()
		}
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if e := root.Eval.Execute(run); assert.NoError(e) {
				lines := assert.Lines()
				assert.Equal(sliceOf.String("hello", "there", "world"), lines)
			}
		}
	}
}

// TestSingleLine verifies the ability of print to join text.
// It complements TestMultiLines
func (assert *CoreSuite) TestSingleLines() {
	var root struct {
		Eval rt.Execute
	}
	if c, ok := assert.ops.NewBuilder(&root); ok {
		if c.Cmd("print span").Begin() {
			if c.Cmds().Begin() {
				if c.Cmd("for each text").Begin() {
					c.Param("in").Val(sliceOf.String("hello", "there", "world"))
					if c.Param("go").Cmds().Begin() {
						c.Cmd("print text", c.Cmd("get", "@", "text")).End()
					}
					c.End()
				}
				c.End()
			}
			c.End()
		}
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if e := root.Eval.Execute(run); assert.NoError(e) {
				lines := assert.Lines()
				assert.Equal("hello there world", lines[0])
			}
		}
	}
}

// TestLineIndex verifies the loop index property.
func (assert *CoreSuite) TestLineIndex() {
	var root struct {
		Eval rt.Execute
	}
	if c, ok := assert.ops.NewBuilder(&root); ok {
		if c.Cmd("for each text").Begin() {
			c.Param("in").Val(sliceOf.String("one", "two", "three"))
			if c.Param("go").Cmds().Begin() {
				if c.Cmd("print num").Begin() {
					c.Cmd("get", "@", "index").End()
				}
				c.End()
			}
			c.End()
		}
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if e := root.Eval.Execute(run); assert.NoError(e) {
				lines := assert.Lines()
				assert.Equal(sliceOf.String("1", "2", "3"), lines)
			}
		}
	}
}

// TestLineEndings verifies loop first and last properties.
func (assert *CoreSuite) TestLineEndings() {
	var root struct {
		Eval rt.Execute
	}
	if c, ok := assert.ops.NewBuilder(&root); ok {
		if c.Cmd("for each text").Begin() {
			c.Param("in").Val(sliceOf.String("one", "two", "three"))
			if c.Param("go").Cmds().Begin() {
				if c.Cmd("print text").Begin() {
					if c.Cmd("choose text").Begin() {
						c.Param("if").Cmd("get", "@", "last")
						c.Param("true").Val("last")
						if c.Param("false").Cmd("choose text").Begin() {
							c.Param("if").Cmd("get", "@", "first")
							c.Param("true").Val("first")
							c.Param("false").Cmd("get", "@", "text")
							c.End()
						}
						c.End()
					}
					c.End()
				}
				c.End()
			}
			c.End()
		}
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if e := root.Eval.Execute(run); assert.NoError(e) {
				lines := assert.Lines()
				assert.Equal(sliceOf.String("first", "two", "last"), lines)
			}
		}
	}
}
