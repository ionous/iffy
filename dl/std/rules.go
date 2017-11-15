package std

import (
	"github.com/ionous/iffy/dl/std/group"
	"github.com/ionous/iffy/spec"
)

// FIX: this has to go into the std library
func Rules(c spec.Block) {
	PrintNameRules(c)
	PrintObjectRules(c)
	group.GroupRules(c)
	storyRules(c)
}

func PrintNameRules(c spec.Block) {
	// print the class name if all else fails
	if c.Cmd("run rule", "print name").Begin() {
		if c.Param("decide").Begin() {
			c.Cmd("say", c.Cmd("class name", c.Cmd("get", "@", "target")))
			c.End()
		}
		c.End()
	}
	// prefer the object name, so long as it was specified by the user.
	if c.Cmd("run rule", "print name").Begin() {
		// detect if "unnamed": # is used only for system names, never author names.
		c.Param("if").Cmd("is not", c.Cmd("includes", c.Cmd("get", c.Cmd("get", "@", "target"), "name"), "#"))
		if c.Param("decide").Begin() {
			c.Cmd("say", c.Cmd("get", c.Cmd("get", "@", "target"), "name"))
			c.End()
		}
		c.End()
	}
	// perfer the printed name above all else
	if c.Cmd("run rule", "print name").Begin() {
		c.Param("if").Cmd("is not", c.Cmd("is empty", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name")))
		if c.Param("decide").Begin() {
			c.Cmd("say", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name"))
			c.End()
		}
		c.End()
	}
	//
	if c.Cmd("run rule", "print plural name").Begin() {
		if c.Param("decide").Begin() {
			if c.Cmd("say").Begin() {
				if c.Cmd("pluralize").Begin() {
					if c.Cmd("buffer").Begin() {
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
		if c.Param("decide").Begin() {
			c.Cmd("say", c.Cmd("get", c.Cmd("get", "@", "target"), "printed plural name"))
			c.End()
		}
		c.End()
	}
	//
	if c.Cmd("run rule", "print several").Begin() {
		if c.Param("decide").Begin() {
			if c.Cmd("print span").Begin() {
				c.Cmd("print num word", c.Cmd("get", "@", "group size"))
				c.Cmd("say", "other")
				if c.Cmd("choose", c.Cmd("compare num", c.Cmd("get", "@", "group size"), c.Cmd("greater than"), 1)).Begin() {
					if c.Param("true").Begin() {
						c.Cmd("determine", c.Cmd("print plural name", c.Cmd("get", "@", "target")))
						c.End()
					}
					if c.Param("false").Begin() {
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
}

func PrintObjectRules(c spec.Block) {
	// print the name and summary if all else fails
	if c.Cmd("run rule", "print object").Begin() {
		if c.Param("decide").Begin() {
			c.Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "target")))
			if c.Cmd("print bracket").Begin() {
				c.Cmd("determine", c.Cmd("print summary", c.Cmd("get", "@", "target")))
				c.End()
			}
			c.End()
		}
		c.End()
	}
	if c.Cmd("run rule", "print summary").Begin() {
		if c.Param("if").Cmd("all true").Begin() {
			c.Cmd("is class", c.Cmd("get", "@", "target"), "container")
			c.Cmd("get", c.Cmd("get", "@", "target"), "closed")
			c.End()
		}
		if c.Param("decide").Begin() {
			c.Cmd("say", "closed")
			c.End()
		}
		c.End()
	}
	// is it better to have multiple patterns, or just one?
	if c.Cmd("run rule", "print summary").Begin() {
		c.Param("if").Cmd("is class", c.Cmd("get", "@", "target"), "container")
		if c.Param("decide").Begin() {
			if c.Cmd("choose", c.Cmd("get", c.Cmd("get", "@", "target"), "closed")).Begin() {
				if c.Param("true").Begin() {
					c.Cmd("say", "closed")
					c.End()
				}
				if c.Param("false").Begin() {
					if c.Cmd("choose", c.Cmd("relation empty", "locale", c.Cmd("get", "@", "target"))).Begin() {
						if c.Param("true").Begin() {
							c.Cmd("say", "open but empty")
							c.End()
						}
						if c.Param("false").Begin() {
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
		if c.Param("decide").Begin() {
			if c.Cmd("print objects").Begin() {
				c.Param("objects").Cmd("related list", "locale", c.Cmd("get", "@", "target"))
				// transfer our print content settings to print objects
				c.Param("header").Cmd("get", "@", "header")
				c.Param("articles").Cmd("get", "@", "articles")
				c.Param("tersely").Cmd("get", "@", "tersely")
				// and handle our fairly magical else
				if c.Param("else").Begin() {
					c.Cmd("say", "empty")
					c.End()
				}
				c.End()
			}
			c.End()
		}
		c.End()
	}
}
