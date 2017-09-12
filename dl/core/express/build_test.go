package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"testing"
)

// TestBuild to ensure that commands are transformed into their appropriate command trees when passing through the builder.
func xTestBuild(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)

	type TestScore struct {
		rt.NumberEval
	}

	unique.PanicTypes(cmds,
		(*TestScore)(nil))
	unique.PanicBlocks(cmds,
		(*std.Commands)(nil), // for renderst
		(*core.Commands)(nil))

	t.Run("property", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, Xform{cmds: cmds})
		c.Cmd("test score", "{val}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{
				&core.GetAt{Prop: "val"},
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
			testEqual(t, &TestScore{
				&core.Get{&core.Global{"Story"}, "score"},
			}, root.NumberEval)
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
				&TestScore{
					&core.Get{&core.Global{"Story"}, "score"},
				}}, root.NumberEval)
		}
	})
}
