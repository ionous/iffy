package spec

import (
	"github.com/ionous/sliceOf"
	"testing"
)

func TestSpec(t *testing.T) {
	cmds := NewCommandBuilder()
	if c := cmds.NewArray(); c.Cmds {
		if c := c.Cmd("unit"); c.Args {
			if c := c.Param("trials").Array(); c.Cmds {
				// cycles
				if c := c.Cmd("match output", sliceOf.String("a", "b", "c", "d")); c.Args {
					if c := c.Cmd("for each num", sliceOf.Float(1, 2, 3, 4)); c.Args {
						if c := c.Cmd("print text"); c.Args {
							c.Cmd("cycle", sliceOf.String("a", "b", "c"))
						}
					}
				}
				// stopping
				if c := c.Cmd("match output", sliceOf.String("a", "b", "c", "c")); c.Args {
					if c := c.Cmd("for each num", sliceOf.Float(1, 2, 3, 4)); c.Args {
						if c := c.Cmd("print text"); c.Args {
							c.Cmd("stopping", sliceOf.String("a", "b", "c"))
						}
					}
				}
			}
		}
	}
	PrintSpec(cmds.Root())
}
