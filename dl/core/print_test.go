package core_test

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/sliceOf"
)

// TestPrintSpacing verifies that print line uses acts like the span printer:
// adding spaces between words as needed.
func (t *CoreSuite) TestPrintSpacing() {
	var root struct {
		Eval rt.Execute
	}
	if c := t.ops.Build(&root); c.Args {
		if c := c.Cmd("print line").Array(); c.Cmds {
			c.Cmd("print text", "hello")
			c.Cmd("print text", "there,")
			c.Cmd("print text", "world.")
		}
	}
	if e := root.Eval.Execute(t.run); t.NoError(e) {
		lines := t.Lines()
		t.Equal("hello there, world.", lines[0], "Note the space after the comma, and the lack of space after the period.")
	}
}

// TestMultiLines verifies that iffy printing works similar to sashimi printing.
// In sashimi, the default printer made every print a new line, we should do the same. This test complements TestSingleLines.
func (t *CoreSuite) TestMultiLines() {
	var root struct {
		Eval rt.Execute
	}
	if c := t.ops.Build(&root); c.Args {
		if c := c.Cmd("for each text"); c.Args {
			c.Param("in").Value(sliceOf.String("hello", "there", "world"))
			if c := c.Param("go").Array(); c.Cmds {
				c.Cmd("print text").Cmd("get", "@", "text")
			}
		}
	}
	if e := root.Eval.Execute(t.run); t.NoError(e) {
		lines := t.Lines()
		t.Equal(sliceOf.String("hello", "there", "world"), lines)
	}
}

// TestSingleLine verifies the ability of print line to join text.
// It complements TestMultiLines
func (t *CoreSuite) TestSingleLines() {
	var root struct {
		Eval rt.Execute
	}
	if c := t.ops.Build(&root); c.Args {
		if c := c.Cmd("print line").Array(); c.Cmds {
			if c := c.Cmd("for each text"); c.Args {
				c.Param("in").Value(sliceOf.String("hello", "there", "world"))
				if c := c.Param("go").Array(); c.Cmds {
					if c := c.Cmd("print text"); c.Args {
						c.Cmd("get", "@", "text")
					}
				}
			}
		}
	}
	if e := root.Eval.Execute(t.run); t.NoError(e) {
		lines := t.Lines()
		t.Equal("hello there world", lines[0])
	}
}

// TestLineIndex verifies the loop index property.
func (t *CoreSuite) TestLineIndex() {
	var root struct {
		Eval rt.Execute
	}
	if c := t.ops.Build(&root); c.Args {
		if c := c.Cmd("for each text"); c.Args {
			c.Param("in").Value(sliceOf.String("one", "two", "three"))
			c.Param("go").Array().Cmd("print num").Cmd("get", "@", "index")
		}
	}
	if e := root.Eval.Execute(t.run); t.NoError(e) {
		lines := t.Lines()
		t.Equal(sliceOf.String("1", "2", "3"), lines)
	}
}

// TestLineEndings verifies loop first and last properties.
func (t *CoreSuite) TestLineEndings() {
	var root struct {
		Eval rt.Execute
	}
	if c := t.ops.Build(&root); c.Args {
		if c := c.Cmd("for each text"); c.Args {
			c.Param("in").Value(sliceOf.String("one", "two", "three"))
			if c := c.Param("go").Array().Cmd("print text"); c.Args {
				if c := c.Cmd("choose text"); c.Args {
					c.Param("if").Cmd("get", "@", "last")
					c.Param("true").Value("last")
					if c := c.Param("false").Cmd("choose text"); c.Args {
						c.Param("if").Cmd("get", "@", "first")
						c.Param("true").Value("first")
						c.Param("false").Cmd("get", "@", "text")
					}
				}
			}
		}
	}
	if e := root.Eval.Execute(t.run); t.NoError(e) {
		lines := t.Lines()
		t.Equal(sliceOf.String("first", "two", "last"), lines)
	}
}
