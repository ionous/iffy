package std_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
	"testing"
)

func TestArticles(t *testing.T) {
	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOps(classes)                   // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types

	unique.PanicBlocks(cmds,
		(*std.Commands)(nil),
		(*core.Commands)(nil),
		(*rule.Commands)(nil),
	)

	unique.PanicBlocks(classes,
		(*std.Classes)(nil))

	unique.PanicBlocks(patterns,
		(*std.Patterns)(nil))

	objects := obj.NewObjects()
	unique.PanicValues(objects,
		&std.Kind{Name: "lamp-post"},
		&std.Kind{Name: "soldiers", IndefiniteArticle: "some"},
		&std.Kind{Name: "trevor", CommonProper: std.ProperNamed},
	)

	rules, e := rule.Master(cmds, core.Xform{}, patterns, std.PrintNameRules)
	if e != nil {
		t.Fatal(e)
	}
	run := rtm.New(classes).Objects(objects).Rules(rules).Rtm()

	match := func(t *testing.T, expected string, fn func(spec.Block)) {
		var lines printer.Lines
		var root struct{ rt.Execute }
		run := rt.Writer(run, &lines)
		if c, ok := cmds.NewXBuilder(&root, core.Xform{}); ok {
			if e := c.Build(fn); e != nil {
				t.Fatal(e)
			} else if e := root.Execute.Execute(run); e != nil {
				t.Fatal(e)
			} else {
				l := lines.Lines()
				if d := pretty.Diff(sliceOf.String(expected), l); len(d) > 0 {
					t.Fatal(d)
				}
			}
		}
	}

	// lower a/n
	t.Run("A trailing lamp post", func(t *testing.T) {
		match(t, "You can only just make out a lamp-post.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", "You can only just make out")
						c.Cmd("say", c.Cmd("lower a/n", "lamp post"))
						c.Cmd("say", ".")
						c.End()
					}
					c.End()
				}
			})
	})

	t.Run("A trailing trevor", func(t *testing.T) {
		match(t, "You can only just make out Trevor.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", "You can only just make out")
						c.Cmd("say", c.Cmd("lower a/n", "trevor"))
						c.Cmd("say", ".")
						c.End()
					}
					c.End()
				}
			})
	})

	t.Run("Trailing some soldiers", func(t *testing.T) {
		match(t, "You can only just make out some soldiers.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", "You can only just make out")
						c.Cmd("say", c.Cmd("lower a/n", "soldiers"))
						c.Cmd("say", ".")
						c.End()
					}
					c.End()
				}
			})
	})

	// upper a/n
	t.Run("A leading lamp post", func(t *testing.T) {
		match(t, "A lamp-post can be made out in the mist.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", c.Cmd("upper a/n", "lamp post"))
						c.Cmd("say", "can be made out in the mist.")
						c.End()
					}
					c.End()
				}
			})
	})

	t.Run("A leading trevor", func(t *testing.T) {
		match(t, "Trevor can be made out in the mist.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", c.Cmd("upper a/n", "trevor"))
						c.Cmd("say", "can be made out in the mist.")
						c.End()
					}
					c.End()
				}
			})
	})

	t.Run("Some leading soldiers", func(t *testing.T) {
		match(t, "Some soldiers can be made out in the mist.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", c.Cmd("upper a/n", "soldiers"))
						c.Cmd("say", "can be made out in the mist.")
						c.End()
					}
					c.End()
				}
			})
	})

	// lower-the
	t.Run("The trailing lamp post", func(t *testing.T) {
		match(t, "You can only just make out the lamp-post.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", "You can only just make out")
						c.Cmd("say", c.Cmd("lower the", "lamp post"))
						c.Cmd("say", ".")
						c.End()
					}
					c.End()
				}
			})
	})

	t.Run("The trailing trevor", func(t *testing.T) {
		match(t, "You can only just make out Trevor.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", "You can only just make out")
						c.Cmd("say", c.Cmd("lower the", "trevor"))
						c.Cmd("say", ".")
						c.End()
					}
					c.End()
				}
			})
	})

	t.Run("The trailing soldiers", func(t *testing.T) {
		match(t, "You can only just make out the soldiers.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", "You can only just make out")
						c.Cmd("say", c.Cmd("lower the", "soldiers"))
						c.Cmd("say", ".")
						c.End()
					}
					c.End()
				}
			})
	})

	// uppe the
	t.Run("The leading lamp post", func(t *testing.T) {
		match(t, "The lamp-post may be a trick of the mist.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", c.Cmd("upper the", "lamp post"))
						c.Cmd("say", "may be a trick of the mist.")
						c.End()
					}
					c.End()
				}
			})
	})

	t.Run("The leading trevor", func(t *testing.T) {
		match(t, "Trevor may be a trick of the mist.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", c.Cmd("upper the", "trevor"))
						c.Cmd("say", "may be a trick of the mist.")
						c.End()
					}
					c.End()
				}
			})
	})

	t.Run("The leading soldiers", func(t *testing.T) {
		match(t, "The soldiers may be a trick of the mist.",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("say", c.Cmd("upper the", "soldiers"))
						c.Cmd("say", "may be a trick of the mist.")
						c.End()
					}
					c.End()
				}
			})
	})

	// FIX: should really be separate -- in a "text" test.
	t.Run("Pluralize", func(t *testing.T) {
		match(t, "lamps",
			func(c spec.Block) {
				if c.Cmd("print span").Begin() {
					c.Cmds(c.Cmd("say", c.Cmd("pluralize", "lamp")))
					c.End()
				}
			})
	})

}
