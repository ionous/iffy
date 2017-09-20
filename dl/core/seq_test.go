package core_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/tests"
	"github.com/ionous/sliceOf"
	"strings"
	"testing"
)

func TestSequences(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil))

	unique.PanicBlocks(classes,
		(*core.Classes)(nil),
	)

	t.Run("cycle text", func(t *testing.T) {
		expect := sliceOf.String("a", "b", "c", "a", "b", "c", "a")
		//
		var n tests.Execute
		var gen obj.Registry
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", len(expect))
			if c.Param("go").Cmds().Begin() {
				c.Cmd("say", c.Cmd("cycle text",
					gen.NewName("cycle counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if run, e := rtm.New(classes).Objects(gen).Rtm(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLines(run, expect); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("shuffle text", func(t *testing.T) {
		var n tests.Execute
		var gen obj.Registry
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", 9)
			if c.Param("go").Cmds().Begin() {
				c.Cmd("say", c.Cmd("shuffle text",
					gen.NewName("shuffle counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()

			if e := c.Build(); e != nil {
				t.Fatal(e)
			} else if run, e := rtm.New(classes).Objects(gen).Rtm(); e != nil {
				t.Fatal(e)
			} else if lines, e := n.GetLines(run); e != nil {
				t.Fatal(e)
			} else if len(lines) != 9 {
				t.Fatal("too few lines", lines)
			} else {
				// ensure every letter appears 3 times
				counter := map[string]int{}
				for _, l := range lines {
					counter[l]++
				}
				for k, v := range counter {
					if v != 3 {
						t.Fatal(k, "should appear equal times")
					}
				}
				// ensure shuffle changes from run to run
				// ex. c,a,b; c,b,a; b,c,a
				c1 := strings.Join(lines[0:3], ",")
				c2 := strings.Join(lines[3:6], ",")
				c3 := strings.Join(lines[6:9], ",")
				if c1 == c2 {
					t.Fatal("1 2 not shuffled", c1, c2, c3)
				} else if c2 == c3 {
					t.Fatal("2 3 not shuffled", c1, c2, c3)
				}
			}
		}
	})
	t.Run("stopping", func(t *testing.T) {
		expect := sliceOf.String("a", "b", "c", "c", "c", "c", "c")
		//
		var n tests.Execute
		var gen obj.Registry
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", len(expect))
			if c.Param("go").Cmds().Begin() {
				c.Cmd("say", c.Cmd("stopping text",
					gen.NewName("stopping counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
			//
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if run, e := rtm.New(classes).Objects(gen).Rtm(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLines(run, expect); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("once", func(t *testing.T) {
		expect := sliceOf.String("a")
		//
		var n tests.Execute
		var gen obj.Registry
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", 5)
			if c.Param("go").Cmds().Begin() {
				c.Cmd("say", c.Cmd("stopping text",
					gen.NewName("stopping counter"),
					sliceOf.String("a"),
				))
				c.End()
			}
			c.End()
			//
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if run, e := rtm.New(classes).Objects(gen).Rtm(); e != nil {
			t.Fatal(e)
		} else if e := n.MatchLines(run, expect); e != nil {
			t.Fatal(e)
		}
	})

}
