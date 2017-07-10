package core_test

import (
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
)

func (assert *CoreSuite) TestSeqCycle() {
	assert.matchLines(sliceOf.String("a", "b", "c", "a", "b", "c", "a"),
		func(c *ops.Builder) {
			if c.Cmd("for each num").Begin() {
				c.Param("in").Cmd("range", 7)
				if c.Param("go").Cmds().Begin() {
					c.Cmd("print text", c.Cmd("cycle text",
						assert.unique.Id("cycle counter"),
						sliceOf.String("a", "b", "c"),
					))
					c.End()
				}
				c.End()
			}
		})
}

func (assert *CoreSuite) TestSeqShuffle() {
	assert.matchFunc(func(c *ops.Builder) {
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", 9)
			if c.Param("go").Cmds().Begin() {
				c.Cmd("print text", c.Cmd("shuffle text",
					assert.unique.Id("shuffle counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
		}
	}, func(lines []string) {
		if assert.Len(lines, 9) {
			counter := map[string]int{}
			for _, l := range lines {
				counter[l]++
			}
			for k, v := range counter {
				if !assert.Equal(3, v, k+" should appear equal times") {
					break
				}
			}
			c1 := lines[0:3]
			c2 := lines[3:6]
			c3 := lines[6:9]
			assert.NotEqual(c1, c2)
			assert.NotEqual(c2, c3)
		}
	})
}

func (assert *CoreSuite) TestSeqStopping() {
	assert.matchLines(
		sliceOf.String("a", "b", "c", "c", "c", "c", "c"),
		func(c *ops.Builder) {
			if c.Cmd("for each num").Begin() {
				c.Param("in").Cmd("range", 7)
				if c.Param("go").Cmds().Begin() {
					c.Cmd("print text", c.Cmd("stopping text",
						assert.unique.Id("stopping counter"),
						sliceOf.String("a", "b", "c"),
					))
					c.End()
				}
				c.End()
			}
		})
}
