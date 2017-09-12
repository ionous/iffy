package std_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/core/express"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/pat/rule"
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
	"testing"
)

func TestStory(t *testing.T) {
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
		&std.Story{Name: "story"},
		&std.Room{Kind: std.Kind{Name: "room"}},
		&std.Pawn{"pawn", ident.IdOf("me")},
		&std.Actor{std.Thing{Kind: std.Kind{Name: "me"}}},
	)
	xform := express.MakeXform(cmds)
	rules, e := rule.Master(cmds, xform, patterns, std.Rules)
	if e != nil {
		t.Fatal(e)
	}

	relations := rel.NewRelations()
	pc := locate.Locale{index.NewTable(index.OneToMany)}
	relations.AddTable("locale", pc.Table)

	run := rtm.New(classes).Objects(objects).Relations(relations).Rules(rules).Rtm()

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
		var root struct{ rt.Execute }
		if c, ok := cmds.NewXBuilder(&root, xform); ok {
			if e := c.Build(fn); e != nil {
				t.Fatal(e)
			} else {
				var lines printer.Lines
				run := rt.Writer(run, &lines)
				if e := root.Execute.Execute(run); e != nil {
					t.Fatal(e)
				} else {
					l := lines.Lines()
					if d := pretty.Diff(sliceOf.String(expected), l); len(d) > 0 {
						t.Log("expected", expected)
						t.Log("got", l)
						t.Fatal(d)
					}
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

}
