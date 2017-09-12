package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"testing"
)

// TestBuild tp ensure we can use templates for eval properties other than text.
func TestBuild(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	patterns := unique.NewStack(cmds.ShadowTypes)

	type TestScore struct {
		rt.NumberEval
	}

	unique.PanicTypes(cmds,
		(*TestScore)(nil))
	unique.PanicBlocks(cmds,
		(*std.Commands)(nil),
		(*core.Commands)(nil),
		(*rule.RuntimeCmds)(nil),
		(*Commands)(nil))
	unique.PanicBlocks(patterns,
		(*std.Patterns)(nil))

	t.Run("property", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, Xform{cmds: cmds})
		c.Cmd("test score", "{val}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{
				&GetAt{Prop: "val"},
			}, root.NumberEval)
		}
	})
	//
	t.Run("global", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, Xform{cmds: cmds})
		c.Cmd("test score", "{Story.score}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{&Render{core.Get{
				&core.Global{"Story"}, "score"},
			}}, root.NumberEval)
		}
	})
	//
	t.Run("run", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, Xform{cmds: cmds})
		c.Cmd("test score", "{go testScore Story.score}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{
				&TestScore{&Render{core.Get{
					&core.Global{"Story"}, "score"},
				}}}, root.NumberEval)
		}
	})
	//
	t.Run("determine", func(t *testing.T) {
		var root struct{ rt.Execute }
		c := cmds.NewBuilder(&root, Xform{cmds: cmds})
		c.Val("{go determine printName Story}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		}
		// hard to compare a shadow class
	})
}
