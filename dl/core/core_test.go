package core_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/tests"
	"testing"
)

func TestCore(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	unique.PanicBlocks(cmds, (*core.Commands)(nil))
	run, e := rtm.New(classes).Rtm()
	if e != nil {
		t.Fatal(e)
	}
	//
	t.Run("text shortcut", func(t *testing.T) {
		var n tests.Text
		c := cmds.NewBuilder(&n, core.Xform{})
		c.Val("some text")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else if e := n.Match(run, "some text"); e != nil {
			t.Fatal(e)
		}
	})
	//
	t.Run("all true", func(t *testing.T) {
		test := func(a, b, res bool) {
			var n tests.Bool
			c := cmds.NewBuilder(&n, core.Xform{})

			c.Cmd("all true", c.Cmds(
				c.Cmd("bool", a),
				c.Cmd("bool", b)))
			//
			if e := c.Build(); e != nil {
				t.Fatal(e)
			} else if e := n.Match(run, res); e != nil {
				t.Fatal(e)
			}
		}
		test(true, false, false)
		test(true, true, true)
		test(false, false, false)
	})
	// ensure AnyTrue operates on boolean literals as "or".
	t.Run("any true", func(t *testing.T) {
		test := func(a, b, res bool) {
			var n tests.Bool
			c := cmds.NewBuilder(&n, core.Xform{})
			c.Cmd("any true", c.Cmds(
				c.Cmd("bool", a),
				c.Cmd("bool", b)))
			//
			if e := c.Build(); e != nil {
				t.Fatal(e)
			} else if e := n.Match(run, res); e != nil {
				t.Fatal(e)
			}
		}
		test(true, false, true)
		test(true, true, true)
		test(false, false, false)
	})
	t.Run("compare numbers", func(t *testing.T) {
		test := func(a float64, op string, b float64, res bool) {
			var n tests.Bool
			c := cmds.NewBuilder(&n, core.Xform{})
			if c.Cmd("compare num").Begin() {
				c.Val(a).Cmd(op).Val(b)
				c.End()
			}
			if e := c.Build(); e != nil {
				t.Fatal(e)
			} else if e := n.Match(run, res); e != nil {
				t.Fatal(e)
			}
		}
		test(10, "greater than", 1, true)
		test(1, "greater than", 10, false)
		test(8, "greater than", 8, false)
		//
		test(10, "lesser than", 1, false)
		test(1, "lesser than", 10, true)
		test(8, "lesser than", 8, false)
		//
		test(10, "not equal to", 1, true)
		test(1, "not equal to", 10, true)
		test(8, "not equal to", 8, false)
		//
		test(10, "equal to", 1, false)
		test(1, "equal to", 10, false)
		test(8, "equal to", 8, true)
	})
	t.Run("compare text", func(t *testing.T) {
		test := func(a string, op string, b string, res bool) {
			var n tests.Bool
			c := cmds.NewBuilder(&n, core.Xform{})
			if c.Cmd("compare text").Begin() {
				c.Val(a).Cmd(op).Val(b)
				c.End()
			}
			if e := c.Build(); e != nil {
				t.Fatal(e)
			} else if e := n.Match(run, res); e != nil {
				t.Fatal(e)
			}
		}
		test("Z", "greater than", "A", true)
		test("A", "greater than", "Z", false)
		//
		test("marzip", "lesser than", "marzipan", true)
		test("marzipan", "lesser than", "marzip", false)
		//
		test("bobby", "not equal to", "phillipa", true)
		test("bobby", "not equal to", "bobby", false)
		//
		test("bobby", "equal to", "bobby", true)
		test("bobby", "equal to", "phillipa", false)
	})
}
