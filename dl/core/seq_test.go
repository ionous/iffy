package core_test

import (
	"github.com/ionous/iffy/ops"
	"github.com/ionous/sliceOf"
)

func (t *CoreSuite) TestCycle() {
	t.match(func(c *ops.Builder) {
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", 4)
			if c.Param("go").Cmds().Begin() {
				c.Cmd("print text", c.Cmd("cycle text",
					t.unique.Id("cycle counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
		}
	}, "a", "b", "c", "a")
}
