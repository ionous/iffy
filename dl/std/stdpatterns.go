package std

import (
	"github.com/ionous/iffy/dl/std/group"
	"github.com/ionous/iffy/spec/ops"
)

// FIX: this has to go into the std library
func StdPatterns(c *ops.Builder) {
	PrintNamePatterns(c)
	PrintObjectPatterns(c)
	group.GroupPatterns(c)
}

func PrintNamePatterns(c *ops.Builder) {
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
					c.Cmd("print text", "other")
					if c.Cmd("choose", c.Cmd("compare num", c.Cmd("get", "@", "group size"), c.Cmd("greater than"), 1)).Begin() {
						if c.Param("true").Cmds().Begin() {
							c.Cmd("determine", c.Cmd("print plural name", c.Cmd("get", "@", "target")))
							c.End()
						}
						if c.Param("false").Cmds().Begin() {
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
}

func PrintObjectPatterns(c *ops.Builder) {
	// print the name and summary if all else fails
	if c.Cmd("run rule", "print object").Begin() {
		if c.Param("decide").Cmds().Begin() {
			c.Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "target")))
			if c.Cmd("print bracket").Begin() {
				c.Cmds(c.Cmd("determine", c.Cmd("print summary", c.Cmd("get", "@", "target"))))
				c.End()
			}
			c.End()
		}
		c.End()
	}
	if c.Cmd("run rule", "print summary").Begin() {
		c.Param("if").Cmd("all true", c.Cmds(
			c.Cmd("is similar class", c.Cmd("get", "@", "target"), "container"),
			c.Cmd("get", c.Cmd("get", "@", "target"), "closed"),
		))
		if c.Param("decide").Cmds().Begin() {
			c.Cmd("print text", "closed")
			c.End()
		}
		c.End()
	}
	// is it better to have multiple patterns, or just one?
	if c.Cmd("run rule", "print summary").Begin() {
		c.Param("if").Cmd("is similar class", c.Cmd("get", "@", "target"), "container")
		if c.Param("decide").Cmds().Begin() {
			if c.Cmd("choose", c.Cmd("get", c.Cmd("get", "@", "target"), "closed")).Begin() {
				c.Param("true").Cmds(c.Cmd("print text", "closed"))
				if c.Param("false").Cmds().Begin() {
					if c.Cmd("choose", c.Cmd("relation empty", "locale", c.Cmd("get", "@", "target"))).Begin() {
						c.Param("true").Cmds(c.Cmd("print text", "open but empty"))
						if c.Param("false").Cmds().Begin() {
							c.Cmd("determine",
								c.Cmd("print content",
									c.Cmd("get", "@", "target"),
									c.Param("header").Val("in which is"),
									c.Param("articles").Val(true),
								))
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
	// keeping the pattern itself bare of parens, etc.
	// that way the caller has some control over how its printed.
	if c.Cmd("run rule", "print content").Begin() {
		if c.Param("decide").Cmds().Begin() {
			if c.Cmd("print objects").Begin() {
				c.Param("objects").Cmd("related list", "locale", c.Cmd("get", "@", "target"))
				// transfer our print content settings to print objects
				c.Param("header").Cmd("get", "@", "header")
				c.Param("articles").Cmd("get", "@", "articles")
				c.Param("tersely").Cmd("get", "@", "tersely")
				// and handle our fairly magical else
				if c.Param("else").Cmds().Begin() {
					c.Cmd("print text", "empty")
					c.End()
				}
				c.End()
			}
			c.End()
		}
		c.End()
	}
}
