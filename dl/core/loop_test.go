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

func TestLoop(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	unique.PanicBlocks(cmds, (*core.Commands)(nil))
	unique.PanicBlocks(classes, (*core.Classes)(nil))
	run, e := rtm.New(classes).Rtm()
	if e != nil {
		t.Fatal(e)
	}

	// verifies the loop index property.
	t.Run("index", func(t *testing.T) {
		var n tests.Execute
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each text").Begin() {
			c.Param("in").Val(sliceOf.String("one", "two", "three"))
			if c.Param("go").Begin() {
				if c.Cmd("print num").Begin() {
					c.Cmd("get", "@", "index").End()
				}
				c.End()
			}
			c.End()
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLines(run, sliceOf.String("1", "2", "3")); e != nil {
			t.Fatal(e)
		}
	})

	// verifies loop first and last properties.
	t.Run("endings", func(t *testing.T) {
		var n tests.Execute
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each text").Begin() {
			c.Param("in").Val(sliceOf.String("one", "two", "three"))
			if c.Param("go").Begin() {
				if c.Cmd("say").Begin() {
					if c.Cmd("choose text").Begin() {
						c.Param("if").Cmd("get", "@", "last")
						c.Param("true").Val("last")
						if c.Param("false").Cmd("choose text").Begin() {
							c.Param("if").Cmd("get", "@", "first")
							c.Param("true").Val("first")
							c.Param("false").Cmd("get", "@", "text")
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
		//
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLines(run, sliceOf.String("first", "two", "last")); e != nil {
			t.Fatal(e)
		}
	})
}
