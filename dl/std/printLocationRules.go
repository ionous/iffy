package std

import "github.com/ionous/iffy/spec"

func printLocationRules(c spec.Block) {
	if c.Cmd("bool rule", "is ceiling").Begin() {
		c.Param("if").Cmd("is class", "object", "room")
		c.Param("decide").Val(true)
		c.End()
	}
	if c.Cmd("bool rule", "is ceiling").Begin() {
		if c.Param("if").Cmd("all true").Begin() {
			c.Cmd("is class", "object", "container")
			c.Cmd("is", "object.closed && !object.transparent")
			c.End()
		}
		c.Param("decide").Val(true)
		c.End()
	}
	if c.Cmd("bool rule", "is notable scenery").Begin() {
		c.Param("decide").Val(false)
		c.End()
	}
	if c.Cmd("bool rule", "is notable enclosure").Begin() {
		c.Param("decide").Val(false)
		c.End()
	}
	if c.Cmd("bool rule", "is unremarkable").Begin() {
		c.Param("if").Cmd("compare obj", "object", "{player}")
		c.Param("decide").Val(true)
		c.End()
	}
	if c.Cmd("bool rule", "is unremarkable").Begin() {
		if c.Param("if").Cmd("any true").Begin() {
			c.Cmd("is").Val("object.handled")
			c.Cmd("is").Val("!object.brief")
			c.End()
		}
		c.Param("decide").Val(true)
		c.End()
	}
	if c.Cmd("run rule", "describe object").Begin() {
		c.Cmd("say", "{printName: object|buffer:}")
		c.End()
	}
	if c.Cmd("run rule", "describe object").Begin() {
		c.Param("if").Cmd("object.brief")
		c.Cmd("say", "{object.brief}")
		c.End()
	}
	if c.Cmd("list objects", "visible parents").Begin() {
		// if c.Cmd("decide").Cmd("list up").Begin() {
		// 	c.Param("source") // parent of object
		// 	//if not this is cieling  parent of object  -- how do i return undecided
		// 	// i could do it in choose if the else is empty-- but what about has parent?
		// 	// maybe something in parents?
		// 	c.Param("next").Cmd(
		// 	c.End()
		// }
		// if c.Cmd("filter", c.Cmd("parents")).Begin() {
		// 	c.Param("accept", "!isCeiling")
		// 	c.End()
		// }

	}
}
