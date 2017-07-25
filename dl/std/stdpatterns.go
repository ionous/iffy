package std

import (
	"github.com/ionous/iffy/dl/std/group"
	"github.com/ionous/iffy/spec/ops"
)

// FIX: this has to go into the std library
func StdPatterns(c *ops.Builder) {
	printPatterns(c)
	group.GroupPatterns(c)
}

func printPatterns(c *ops.Builder) {
	// its a little heavy to do this with patterns, but -- its a good test of the system.
	// print the class name if all else fails
	if c.Cmd("run rule", "print name").Begin() {
		if c.Param("decide").Cmds().Begin() {
			c.Cmd("print text", c.Cmd("class name", c.Cmd("get", "@", "target")))
			c.End()
		}
		c.End()
	}
	// prefer the object name, so long as it was specified by the user.
	if c.Cmd("run rule", "print name").Begin() {
		// detect if "unnamed": # is used only for system names, never author names.
		c.Param("if").Cmd("is not", c.Cmd("includes", c.Cmd("get", c.Cmd("get", "@", "target"), "name"), "#"))
		if c.Param("decide").Cmds().Begin() {
			c.Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "name"))
			c.End()
		}
		c.End()
	}
	// perfer the printed name above all else
	if c.Cmd("run rule", "print name").Begin() {
		c.Param("if").Cmd("is not", c.Cmd("is empty", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name")))
		if c.Param("decide").Cmds().Begin() {
			c.Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name"))
			c.End()
		}
		c.End()
	}
	//
	if c.Cmd("run rule", "print plural name").Begin() {
		if c.Param("decide").Cmds().Begin() {
			if c.Cmd("print text").Begin() {
				if c.Cmd("pluralize").Begin() {
					if c.Cmd("buffer").Begin() {
						if c.Cmds().Begin() {
							c.Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "target")))
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
		c.End()
	}
	if c.Cmd("run rule", "print plural name").Begin() {
		c.Param("if").Cmd("is not", c.Cmd("is empty", c.Cmd("get", c.Cmd("get", "@", "target"), "printed plural name")))
		if c.Param("decide").Cmds().Begin() {
			c.Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "printed plural name"))
			c.End()
		}
		c.End()
	}
	//
	if c.Cmd("run rule", "print several").Begin() {
		if c.Param("decide").Cmds().Begin() {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print num word", c.Cmd("get", "@", "group size"))
					c.Cmd("determine", c.Cmd("print plural name", c.Cmd("get", "@", "target")))
					c.End()
				}
				c.End()
			}
			c.End()
		}
		c.End()
	}
}
