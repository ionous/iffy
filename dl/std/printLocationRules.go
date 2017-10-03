package std

import "github.com/ionous/iffy/spec"

func printLocationRules(c spec.Block) {
	if c.Cmd("bool rule", "is ceiling").Begin() {
		c.Param("if").Cmd("is class", "object", "room")
		c.Param("decide").Val(true)
	}
	if c.Cmd("bool rule", "is ceiling").Begin() {
		if c.Param("if").Cmd("all true").Begin() {
			if c.Cmds().Begin() {
				c.Cmd("is class", "object", "container")
				c.Cmd("is", "{object.closed && !object.transparent}")
				c.End()
			}
			c.End()
		}
		c.Param("decide").Val(true)
		c.End()
	}

}
