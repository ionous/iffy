package std

import (
	"github.com/ionous/iffy/spec/ops"
)

// FIX: this has to go into the std library
func BuildPatterns(c *ops.Builder) {
	namePatterns(c)
	groupPatterns(c)
}

func namePatterns(c *ops.Builder) {
	// its a little heavy to do this with patterns, but -- its a good test of the system.
	// print the class name if all else fails
	if c.Cmd("run rule", "print name").Begin() {
		c.Param("decide").Cmd("print text", c.Cmd("class name", c.Cmd("get", "@", "target")))
		c.End()
	}
	// prefer the object name, so long as it was specified by the user.
	if c.Cmd("run rule", "print name").Begin() {
		// # is used only for system names, not user author names.
		c.Param("if").Cmd("is not", c.Cmd("includes", c.Cmd("get", c.Cmd("get", "@", "target"), "name"), "#"))
		c.Param("decide").Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "name"))
		c.End()
	}
	// perfer the printed name above all else
	if c.Cmd("run rule", "print name").Begin() {
		c.Param("if").Cmd("is not", c.Cmd("is empty", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name")))
		c.Param("decide").Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name"))
		c.End()
	}
	//
	if c.Cmd("run rule", "print plural name").Begin() {
		if c.Param("decide").Cmd("print text").Begin() {
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
	if c.Cmd("run rule", "print plural name").Begin() {
		c.Param("if").Cmd("is not", c.Cmd("is empty", c.Cmd("get", c.Cmd("get", "@", "target"), "printed plural name")))
		c.Param("decide").Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "printed plural name"))
		c.End()
	}
}

func groupPatterns(c *ops.Builder) {

}
