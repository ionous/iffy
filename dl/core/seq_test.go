package core_test

import (
	"github.com/ionous/iffy/ops"
	"github.com/ionous/sliceOf"
)

func (t *CoreSuite) TestSeqCycle() {
	t.match(func(c *ops.Builder) {
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", 7)
			if c.Param("go").Cmds().Begin() {
				c.Cmd("print text", c.Cmd("cycle text",
					t.unique.Id("cycle counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
		}
	}, "a", "b", "c", "a", "b", "c", "a")
}

func (t *CoreSuite) TestSeqShuffle() {
	t.matchFunc(func(c *ops.Builder) {
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", 9)
			if c.Param("go").Cmds().Begin() {
				c.Cmd("print text", c.Cmd("shuffle text",
					t.unique.Id("shuffle counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
		}
	}, func(lines []string) {
		if t.Len(lines, 9) {
			counter := map[string]int{}
			for _, l := range lines {
				counter[l]++
			}
			for k, v := range counter {
				if !t.Equal(3, v, k+" should appear equal times") {
					break
				}
			}
			c1 := lines[0:3]
			c2 := lines[3:6]
			c3 := lines[6:9]
			t.NotEqual(c1, c2)
			t.NotEqual(c2, c3)
		}
	})
}

func (t *CoreSuite) TestSeqStopping() {
	t.match(func(c *ops.Builder) {
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", 7)
			if c.Param("go").Cmds().Begin() {
				c.Cmd("print text", c.Cmd("stopping text",
					t.unique.Id("stopping counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
		}
	}, "a", "b", "c", "c", "c", "c", "c")
}
