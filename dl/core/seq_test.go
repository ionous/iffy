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
	gen := unique.NewObjectGenerator()

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil))

	unique.PanicBlocks(gen,
		(*core.Counters)(nil))

	unique.PanicBlocks(classes,
		(*core.Classes)(nil),
		(*core.Counters)(nil),
	)

	panicGen := func() *obj.ObjBuilder {
		objs, e := gen.Generate()
		if e != nil {
			panic(e)
		}
		objects := obj.NewObjects()
		unique.PanicValues(objects, objs...)
		return objects
	}

	t.Run("cycle text", func(t *testing.T) {
		expect := sliceOf.String("a", "b", "c", "a", "b", "c", "a")
		//
		var n tests.Execute
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", len(expect))
			if c.Param("go").Cmds().Begin() {
				c.Cmd("say", c.Cmd("cycle text",
					gen.Id("cycle counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			run := rtm.New(classes).Objects(panicGen()).Rtm()
			if e := n.MatchLines(run, expect); e != nil {
				t.Fatal(e)
			}
		}
	})
	t.Run("shuffle text", func(t *testing.T) {
		var n tests.Execute
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", 9)
			if c.Param("go").Cmds().Begin() {
				c.Cmd("say", c.Cmd("shuffle text",
					gen.Id("shuffle counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()

			if e := c.Build(); e != nil {
				t.Fatal(e)
			} else {
				run := rtm.New(classes).Objects(panicGen()).Rtm()
				if lines, e := n.GetLines(run); e != nil {
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
		}
	})
	t.Run("stopping", func(t *testing.T) {
		expect := sliceOf.String("a", "b", "c", "c", "c", "c", "c")
		//
		var n tests.Execute
		c := cmds.NewBuilder(&n, core.Xform{})
		if c.Cmd("for each num").Begin() {
			c.Param("in").Cmd("range", len(expect))
			if c.Param("go").Cmds().Begin() {
				c.Cmd("say", c.Cmd("stopping text",
					gen.Id("stopping counter"),
					sliceOf.String("a", "b", "c"),
				))
				c.End()
			}
			c.End()
			//
		}
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			run := rtm.New(classes).Objects(panicGen()).Rtm()
			if e := n.MatchLines(run, expect); e != nil {
				t.Fatal(e)
			}
		}
	})
}
