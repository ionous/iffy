package core_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/tests"
	"github.com/ionous/sliceOf"
	"testing"
)

func TestPrint(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	unique.PanicBlocks(cmds, (*core.Commands)(nil))
	run, e := rtm.New(classes).Rtm()
	if e != nil {
		t.Fatal(e)
	}
	//
	t.Run("spacing", func(t *testing.T) {
		var n tests.Execute
		c := cmds.NewBuilder(&n, ops.Transformer(core.Transform))
		if c.Cmd("print span").Begin() {
			c.Cmd("say", "hello")
			c.Cmd("say", "there,")
			c.Cmd("say", "world.")
			c.End()
		}
		//
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLine(run, "hello there, world."); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("numbers", func(t *testing.T) {
		var n tests.Execute
		c := cmds.NewBuilder(&n, ops.Transformer(core.Transform))
		if c.Cmd("print span").Begin() {
			c.Cmd("print num", 213)
			c.Cmd("print num word", 213)
			c.End()
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLine(run, "213 two hundred thirteen"); e != nil {
			t.Fatal(e)
		}
	})

	// In sashimi, the default printer made every print a new line, we should do the same. This test complements TestSingleLines.
	t.Run("multi", func(t *testing.T) {
		var n tests.Execute
		c := cmds.NewBuilder(&n, ops.Transformer(core.Transform))
		if c.Cmd("for each text").Begin() {
			c.Param("in").Val(sliceOf.String("hello", "there", "world"))
			if c.Param("go").Begin() {
				c.Cmd("say", c.Cmd("get", "@", "text"))
				c.End()
			}
			c.End()
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLines(run, sliceOf.String("hello", "there", "world")); e != nil {
			t.Fatal(e)
		}
	})

	// TestSingleLine verifies the ability of print to join text.
	// It complements TestMultiLines
	t.Run("single", func(t *testing.T) {
		var n tests.Execute
		c := cmds.NewBuilder(&n, ops.Transformer(core.Transform))
		if c.Cmd("print span").Begin() {
			if c.Cmd("for each text").Begin() {
				c.Param("in").Val(sliceOf.String("hello", "there", "world"))
				if c.Param("go").Begin() {
					c.Cmd("say", c.Cmd("get", "@", "text")).End()
				}
				c.End()
			}
			c.End()
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLine(run, "hello there world"); e != nil {
			t.Fatal(e)
		}
	})
}
