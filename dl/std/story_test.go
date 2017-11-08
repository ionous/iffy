package std_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/express"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/rel"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
	"strings"
	"testing"
)

func xTestStory(t *testing.T) {
	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOps(classes)                   // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*std.Commands)(nil),
		(*express.Commands)(nil),
		(*rules.Commands)(nil),
	)
	unique.PanicBlocks(classes,
		(*std.Classes)(nil))

	unique.PanicBlocks(patterns,
		(*std.Patterns)(nil))

	var objects obj.Registry
	story := &std.Story{Name: "story"}
	objects.RegisterValues(sliceOf.Interface(
		story,
		&std.Room{Kind: std.Kind{Name: "room"}},
		&std.Pawn{"pawn", ident.IdOf("me")},
		&std.Actor{std.Thing{Kind: std.Kind{Name: "me"}}},
	))
	xform := express.NewTransform(cmds, nil)
	rules, e := rules.Master(cmds, xform, patterns, std.Rules)
	if e != nil {
		t.Fatal(e)
	}

	relations := rel.NewRelations()
	pc := locate.Locale{index.NewTable(index.OneToMany)}
	relations.AddTable("locale", pc.Table)

	run, e := rtm.New(classes).Objects(objects).Relations(relations).Rules(rules).Rtm()
	if e != nil {
		t.Fatal(e)
	}

	Object := func(name string) rt.Object {
		ret, ok := run.GetObject(name)
		if !ok {
			t.Fatal("couldnt find object", name)
		}
		return ret
	}
	if e := pc.SetLocation(Object("room"), locate.Has, Object("me")); e != nil {
		t.Fatal(e)
	}

	match := func(t *testing.T, expected string, fn func(spec.Block)) {
		var root struct{ rt.ExecuteList }
		c := cmds.NewBuilder(&root, xform)
		if e := c.Build(fn); e != nil {
			t.Fatal(e)
		} else {
			// t.Log(pretty.Sprint(root.ExecuteList))
			var lines printer.Lines
			if e := rt.WritersBlock(run, &lines, func() error {
				return root.Execute(run)
			}); e != nil {
				t.Fatal(e)
			} else {
				l := strings.Join(lines.Lines(), "\n")
				if d := pretty.Diff(expected, l); len(d) > 0 {
					t.Log("expected", expected)
					t.Log("got", l)
					t.Fatal(d)
				}
			}
		}
	}
	t.Run("print location", func(t *testing.T) {
		match(t, "room", func(c spec.Block) {
			c.Cmd("determine", c.Cmd("print name", c.Cmd("location of", c.Cmd("player"))))
		})
	})
	t.Run("surroundings", func(t *testing.T) {
		match(t, "room", func(c spec.Block) {
			if c.Cmd("say").Begin() {
				c.Cmd("determine", c.Cmd("player surroundings"))
				c.End()
			}
		})
	})
	t.Run("status left", func(t *testing.T) {
		match(t, "room", func(c spec.Block) {
			c.Cmd("set text", "story", "status left", "{determine playerSurroundings}")
			c.Cmd("say", "{story.statusLeft}")
		})
	})
	t.Run("status right", func(t *testing.T) {
		match(t, "0/0", func(c spec.Block) {
			c.Cmd("set text", "story", "status right", "{score}/{turnCount}")
			c.Cmd("say", "{story.statusRight}")
		})
	})
	t.Run("banner defaults", func(t *testing.T) {
		x := strings.Join(sliceOf.String(
			"Welcome",
			"An interactive fiction",
			"Release 0.0.0 / Iffy 1.0",
		), "\n")
		match(t, x, func(c spec.Block) {
			c.Cmd("determine", c.Cmd("print banner text"))
		})
	})
	t.Run("banner text", func(t *testing.T) {
		story.Title = "Curses"
		story.Author = "An other mouse"
		story.MajorVersion = 1
		story.MinorVersion = 2
		story.PatchVersion = 3
		story.SerialNumber = "YYMMDD"
		x := strings.Join(sliceOf.String(
			"Curses",
			"An interactive fiction by An other mouse",
			"Release 1.2.3 / YYMMDD / Iffy 1.0",
		), "\n")
		match(t, x, func(c spec.Block) {
			c.Cmd("determine", c.Cmd("print banner text"))
		})
	})

}
