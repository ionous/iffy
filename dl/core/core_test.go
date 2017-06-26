package core_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ops"
	"github.com/ionous/iffy/ops/unique"
	"github.com/ionous/iffy/reflector"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

func TestCore(t *testing.T) {
	suite.Run(t, new(CoreSuite))
}

// Regular expression to select test suites specified command-line argument "-run". Regular expression to select the methods of test suites specified command-line argument "-m"
type CoreSuite struct {
	suite.Suite
	ops  *ops.Ops
	test *testing.T

	run   rt.Runtime
	lines rtm.LineWriter

	unique *unique.Objects
}

func (t *CoreSuite) Lines() (ret []string) {
	ret = t.lines.Lines()
	t.lines = rtm.LineWriter{}
	return
}

func (t *CoreSuite) SetupTest() {
	t.ops = ops.NewOps((*core.Commands)(nil))
	t.test = t.T()
	t.unique = unique.NewObjects()
	t.unique.Types.RegisterType((*core.CycleCounter)(nil))
}

func (t *CoreSuite) newRuntime(c *ops.Builder) (ret rt.Runtime, err error) {
	if _, e := c.Build(); e != nil {
		err = e
	} else {

		mm := reflector.NewModelMaker()
		mm.AddClass((*core.NumberCounter)(nil))
		mm.AddClass((*core.TextCounter)(nil))
		mm.AddClass((*core.CycleCounter)(nil))

		if inst, e := t.unique.Generate(); e != nil {
			err = e
		} else {
			if cnt := len(inst); cnt > 0 {
				t.test.Log("creating", cnt, "instances")
				mm.AddInstance(inst...)
			}
			//
			if m, e := mm.MakeModel(); e != nil {
				err = e
			} else {
				run := rtm.NewRtm(m)
				run.PushWriter(&t.lines)
				ret = run
			}
		}
	}
	return
}

func (t *CoreSuite) match(
	build func(c *ops.Builder), expected ...string) {
	var root struct{ Eval rt.Execute }
	if c, ok := t.ops.NewBuilder(&root); ok {
		build(c)
		if run, e := t.newRuntime(c); t.NoError(e) {
			if e := root.Eval.Execute(run); t.NoError(e) {
				lines := t.Lines()
				t.Equal(expected, lines)
			}
		}
	}
}

func (t *CoreSuite) TestShortcuts() {
	var root struct {
		Eval rt.TextEval
	}
	if c, ok := t.ops.NewBuilder(&root); ok {
		c.Val("shortcut")
		if run, e := t.newRuntime(c); t.NoError(e) {
			if res, e := root.Eval.GetText(run); t.NoError(e) {
				t.EqualValues("shortcut", res)
			}
		}
	}
}

// TestAllTrue ensure AllTrue operates on boolean literals as "and".
func (t *CoreSuite) TestAllTrue() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a, b, res bool) {
		if c, ok := t.ops.NewBuilder(&root); ok {
			c.Cmd("all true", c.Cmds(
				c.Cmd("bool", a),
				c.Cmd("bool", b)))
			//
			if run, e := t.newRuntime(c); t.NoError(e) {
				if ok, e := root.Eval.GetBool(run); t.NoError(e) {
					t.EqualValues(res, ok)
				}
			}
		}
	}
	test(true, false, false)
	test(true, true, true)
	test(false, false, false)
}

// TestAnyTrue ensure AnyTrue operates on boolean literals as "or".
func (t *CoreSuite) TestAnyTrue() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a, b, res bool) {
		if c, ok := t.ops.NewBuilder(&root); ok {
			if c.Cmd("any true").Begin() {
				c.Cmds(c.Cmd("bool", a), c.Cmd("bool", b))
				c.End()
			}
			// /
			if run, e := t.newRuntime(c); t.NoError(e) {
				if ok, e := root.Eval.GetBool(run); t.NoError(e) {
					t.EqualValues(res, ok)
				}
			}
		}
	}
	test(true, false, true)
	test(true, true, true)
	test(false, false, false)
}

func (t *CoreSuite) TestCompareNum() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a float64, op string, b float64) {
		if c, ok := t.ops.NewBuilder(&root); ok {
			if c.Cmd("compare num").Begin() {
				c.Val(a).Cmd(op).Val(b)
				c.End()
			}

			if run, e := t.newRuntime(c); t.NoError(e) {
				if ok, e := root.Eval.GetBool(run); t.NoError(e) {
					t.True(ok)
				}
			}
		}
	}
	test(10, "greater than", 1)
	test(1, "lesser than", 10)
	test(10, "not equal to", 1)
	test(10, "equal to", 10)
}

func (t *CoreSuite) TestCompareText() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a, op, b string) {
		if c, ok := t.ops.NewBuilder(&root); ok {
			c.Cmd("compare text", c.Val(a), c.Cmd(op), c.Val(b))
			//
			if run, e := t.newRuntime(c); t.NoError(e) {
				if ok, e := root.Eval.GetBool(run); t.NoError(e) {
					t.True(ok, strings.Join([]string{a, op, b}, " "))
				}
			}
		}
	}
	test("Z", "greater than", "A")
	test("marzip", "lesser than", "marzipan")
	test("bobby", "not equal to", "suzzie")
	test("tyrone", "equal to", "tyrone")
}
