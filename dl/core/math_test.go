package core_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/tests"
	"testing"
)

func TestMath(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	//
	unique.PanicBlocks(cmds,
		(*core.Commands)(nil))
	unique.PanicBlocks(classes,
		(*core.Classes)(nil),
	)
	//
	run, e := rtm.New(classes).Rtm()
	if e != nil {
		t.Fatal(e)
	}
	//
	match := func(t *testing.T, v float64, fn func(spec.Block)) (err error) {
		var n tests.Number
		c := cmds.NewBuilder(&n, core.Xform{})
		if e := c.Build(fn); e != nil {
			err = e
		} else if e := n.Match(run, v); e != nil {
			err = e
		}
		return
	}
	t.Run("Add", func(t *testing.T) {
		if e := match(t, 11, func(c spec.Block) {
			c.Cmd("add", 1, 10)
		}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Sub", func(t *testing.T) {
		if e := match(t, -9, func(c spec.Block) {
			c.Cmd("sub", 1, 10)
		}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Mul", func(t *testing.T) {
		if e := match(t, 200, func(c spec.Block) {
			c.Cmd("mul", 20, 10)
		}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Div", func(t *testing.T) {
		if e := match(t, 2, func(c spec.Block) {
			c.Cmd("div", 20, 10)
		}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Div By Zero", func(t *testing.T) {
		if e := match(t, 200, func(c spec.Block) {
			c.Cmd("div", 20, 10)
		}); e == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("Mod", func(t *testing.T) {
		if e := match(t, 1, func(c spec.Block) {
			c.Cmd("mod", 3, 2)
		}); e != nil {
			t.Fatal(e)
		}
	})
}
