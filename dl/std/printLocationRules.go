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
			c.Cmd("is", "{object.closed && !object.transparent}")
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
		c.Param("if").Cmd("compare obj", "{object}", "{pawn.actor}")
		c.Param("decide").Val(true)
		c.End()
	}
	if c.Cmd("bool rule", "is unremarkable").Begin() {
		if c.Param("if").Cmd("any true").Begin() {
			c.Cmd("is").Val("{object.handled}")
			c.Cmd("is").Val("{!object.brief}")
			c.End()
		}
		c.Param("decide").Val(true)
		c.End()
	}
	// if c.Cmd("run rule", "describe object").Begin() {
	// 	if c.Param("decide").Begin() {
	// 		if c.Cmd("choose").Begin() {
	// 			c.Param("if").Cmd("{object.brief}")
	// 			c.Param("true").Cmd("say", "{object.brief}")
	// 			c.Param("false").Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "target")))
	// 			c.End()
	// 		}
	// 	}
	// 	c.End()
	// }

}
