package std_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/express"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/dl/std"
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
		(*rules.Commands)(nil),
		(*express.Commands)(nil),
	)

	unique.PanicBlocks(classes,
		(*std.Classes)(nil))

	unique.PanicBlocks(patterns,
		(*std.Patterns)(nil))

	var objects obj.Registry
	objects.RegisterValues(sliceOf.Interface(
		&std.Kind{Name: "lamp-post"},
		&std.Kind{Name: "soldiers", IndefiniteArticle: "some"},
		&std.Kind{Name: "trevor", CommonProper: std.ProperNamed},
	))
	xform := express.NewTransform(cmds, nil)
	rules, e := rules.Master(cmds, xform, patterns, std.PrintNameRules)
	if e != nil {
		t.Fatal(e)
	}
	run, e := rtm.New(classes).Objects(objects).Rules(rules).Rtm()
	if e != nil {
		t.Fatal(e)
	}

	match := func(t *testing.T, expected string, fn func(spec.Block)) error {
		var lines printer.Lines
		return rt.WritersBlock(run, &lines, func() (err error) {
			var root struct{ rt.Execute }
			c := cmds.NewBuilder(&root, xform)
			if e := c.Build(fn); e != nil {
				err = e
			} else if e := root.Execute.Execute(run); e != nil {
				err = e
			} else {
				l := lines.Lines()
				t.Log(l)
				if d := pretty.Diff(sliceOf.String(expected), l); len(d) > 0 {
					err = errutil.New("mismatch", d)
				}
			}
			return
		})
	}
	say := func(t *testing.T, expected, say string) {
		if e := match(t, expected, func(c spec.Block) {
			c.Cmd("say", say)
		}); e != nil {
			t.Fatal(e)
		}
	}

	// lower a/n
	t.Run("A trailing lamp post", func(t *testing.T) {
		say(t, "You can only just make out a lamp-post.",
			"You can only just make out {lowerAn: lampPost}.")
		// if c.Cmd("print span").Begin() {
		// 	c.Cmd("say", "You can only just make out")
		// 	c.Cmd("say", c.Cmd("lower a/n", "lamp post"))
		// 	c.Cmd("say", ".")
		// 	c.End()
		// }
	})

	t.Run("A trailing trevor", func(t *testing.T) {
		say(t, "You can only just make out Trevor.",
			"You can only just make out {lowerAn: trevor}.")
		// if c.Cmd("print span").Begin() {
		// 	c.Cmd("say", "You can only just make out")
		// 	c.Cmd("say", c.Cmd("lower a/n", "trevor"))
		// 	c.Cmd("say", ".")
		// 	c.End()
		// }
	})

	t.Run("Trailing some soldiers", func(t *testing.T) {
		say(t, "You can only just make out some soldiers.",
			"You can only just make out {lowerAn: soldiers}.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", "You can only just make out")
		// 		c.Cmd("say", c.Cmd("lower a/n", "soldiers"))
		// 		c.Cmd("say", ".")
		// 		c.End()
		// 	}
		// })
	})

	// upper a/n
	t.Run("A leading lamp post", func(t *testing.T) {
		say(t, "A lamp-post can be made out in the mist.",
			"{upperAn: lampPost} can be made out in the mist.")
		// if c.Cmd("print span").Begin() {
		// 	c.Cmd("say", c.Cmd("upper a/n", "lamp post"))
		// 	c.Cmd("say", "can be made out in the mist.")
		// 	c.End()
		// }
		// })
	})

	t.Run("A leading trevor", func(t *testing.T) {
		say(t, "Trevor can be made out in the mist.",
			"{upperAn: trevor} can be made out in the mist.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", c.Cmd("upper a/n", "trevor"))
		// 		c.Cmd("say", "can be made out in the mist.")
		// 		c.End()
		// 	}
		// })
	})

	t.Run("Some leading soldiers", func(t *testing.T) {
		say(t, "Some soldiers can be made out in the mist.",
			"{upperAn: soldiers} can be made out in the mist.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", c.Cmd("upper a/n", "soldiers"))
		// 		c.Cmd("say", "can be made out in the mist.")
		// 		c.End()
		// 	}
		// })
	})

	// lower-the
	t.Run("The trailing lamp post", func(t *testing.T) {
		say(t, "You can only just make out the lamp-post.",
			"You can only just make out {lowerThe: lampPost}.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", "You can only just make out")
		// 		c.Cmd("say", c.Cmd("lower the", "lamp post"))
		// 		c.Cmd("say", ".")
		// 		c.End()
		// 	}
		// })
	})

	t.Run("The trailing trevor", func(t *testing.T) {
		say(t, "You can only just make out Trevor.",
			"You can only just make out {lowerThe: trevor}.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", "You can only just make out")
		// 		c.Cmd("say", c.Cmd("lower the", "trevor"))
		// 		c.Cmd("say", ".")
		// 		c.End()
		// 	}
		// })
	})

	t.Run("The trailing soldiers", func(t *testing.T) {
		say(t, "You can only just make out the soldiers.",
			"You can only just make out {lowerThe: soldiers}.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", "You can only just make out")
		// 		c.Cmd("say", c.Cmd("lower the", "soldiers"))
		// 		c.Cmd("say", ".")
		// 		c.End()
		// 	}
		// })
	})

	// uppe the
	t.Run("The leading lamp post", func(t *testing.T) {
		say(t, "The lamp-post may be a trick of the mist.",
			"{upperThe: lampPost} may be a trick of the mist.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", c.Cmd("upper the", "lamp post"))
		// 		c.Cmd("say", "may be a trick of the mist.")
		// 		c.End()
		// 	}
		// })
	})

	t.Run("The leading trevor", func(t *testing.T) {
		say(t, "Trevor may be a trick of the mist.",
			"{upperThe: trevor} may be a trick of the mist.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", c.Cmd("upper the", "trevor"))
		// 		c.Cmd("say", "may be a trick of the mist.")
		// 		c.End()
		// 	}
		// })
	})

	t.Run("The leading soldiers", func(t *testing.T) {
		say(t, "The soldiers may be a trick of the mist.",
			"{upperThe: soldiers} may be a trick of the mist.")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", c.Cmd("upper the", "soldiers"))
		// 		c.Cmd("say", "may be a trick of the mist.")
		// 		c.End()
		// 	}
		// })
	})

	// FIX: should really be separate -- in a "text" test.
	t.Run("Pluralize", func(t *testing.T) {
		say(t, "lamps", "{pluralize: 'lamp'}")
		// func(c spec.Block) {
		// 	if c.Cmd("print span").Begin() {
		// 		c.Cmd("say", c.Cmd("pluralize", "lamp"))
		// 		c.End()
		// 	}
		// })
	})
}
