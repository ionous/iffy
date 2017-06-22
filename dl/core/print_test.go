package core_test

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/sliceOf"
)

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

// in the original, the default printer made every print a new line
// we should do the same.
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

// // needs ChooseText and substitution
// func (t *CoreSuite) TestLineEndings() {
// 	var root struct {
// 		Eval rt.Execute
// 	}
// 	if c := t.ops.Build(&root); c.Args {
// 		if c := c.Cmd("for each text"); c.Cmds {
// 			c.Cmd("texts", sliceOf.String("one", "two", "three"))
// 			if c := c.Cmd("print text"); c.Args {
// 				if c := c.Cmd("choose text"); c.Args {
// 					c.Param("if").Cmd("get", "@", "first")
// 					c.Param("true").Value("first")
// 					if c := c.Param("false").Cmd("choose text"); c.Args {
// 						c.Param("if").Cmd("get", "@", "last")
// 						c.Param("true").Value("last")
// 						c.Param("false").Cmd("get", "@", "text")
// 					}
// 				}
// 				c.Cmd("object", "@")
// 				c.Value("num")
// 			}
// 		}
// 	}
// 	if e := root.Eval.Execute(t.run); t.NoError(e) {
// 		lines := t.Lines()
// 		t.Equal(sliceOf.String("first",
// 			"two",
// 			"last"), lines)
// 	}
// }
