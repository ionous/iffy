package core_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ops"
	"github.com/ionous/iffy/reflector"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	. "github.com/ionous/iffy/tests"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

var people = []interface{}{
	&Person{"alice", Listening, Standing},
	&Person{"bob", Listening, Sitting},
	&Person{"carol", Laughing, Sitting},
}

// Regular expression to select test suites specified command-line argument "-run". Regular expression to select the methods of test suites specified command-line argument "-m"
type CoreSuite struct {
	suite.Suite
	ops   *ops.Ops
	run   rt.Runtime
	lines rtm.LineWriter
}

func TestCoreSuite(t *testing.T) {
	suite.Run(t, new(CoreSuite))
}

func (t *CoreSuite) Log(args ...interface{}) {
	t.T().Log(args...)
}
func (t *CoreSuite) Lines() (ret []string) {
	ret = t.lines.Lines()
	t.lines = rtm.LineWriter{}
	return
}

func (t *CoreSuite) SetupTest() {
	errutil.Panic = true
	t.ops = ops.NewOps((*core.Commands)(nil))
	mm := reflector.NewModelMaker()
	mm.AddClass((*core.NumberCounter)(nil))
	mm.AddClass((*core.TextCounter)(nil))
	//
	if m, e := mm.MakeModel(); e != nil {
		panic(e)
	} else {
		t.run = rtm.NewRtm(m)
		t.run.PushWriter(&t.lines)
	}
}

func (t *CoreSuite) TestShortcuts() {
	var root struct {
		Eval rt.TextEval
	}
	if c := t.ops.Build(&root); c.Args {
		c.Value("shortcut")
	}
	if res, e := root.Eval.GetText(t.run); t.NoError(e) {
		t.EqualValues("shortcut", res)
	}
}

// TestAllTrue ensure AllTrue operates on boolean literals as "and".
func (t *CoreSuite) TestAllTrue() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a, b, c bool) {
		if c := t.ops.Build(&root); c.Args {
			if c := c.Cmd("all true").Array(); c.Cmds {
				c.Cmd("bool", a)
				c.Cmd("bool", b)
			}
		}
		if ok, e := root.Eval.GetBool(t.run); t.NoError(e) {
			t.EqualValues(c, ok)
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
	test := func(a, b, c bool) {
		if c := t.ops.Build(&root); c.Args {
			if c := c.Cmd("any true").Array(); c.Cmds {
				c.Cmd("bool", a)
				c.Cmd("bool", b)
			}
		}
		if ok, e := root.Eval.GetBool(t.run); t.NoError(e) {
			t.EqualValues(c, ok)
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
		if c := t.ops.Build(&root); c.Args {
			if c := c.Cmd("compare num"); c.Args {
				c.Cmd("num", a)
				c.Cmd(op)
				c.Cmd("num", b)
			}
		}
		if ok, e := root.Eval.GetBool(t.run); t.NoError(e) {
			t.True(ok)
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
		if c := t.ops.Build(&root); c.Args {
			if c := c.Cmd("compare text"); c.Args {
				c.Cmd("text", a)
				c.Cmd(op)
				c.Cmd("text", b)
			}
		}
		if ok, e := root.Eval.GetBool(t.run); t.NoError(e) {
			t.True(ok, strings.Join([]string{a, op, b}, " "))
		}
	}
	test("Z", "greater than", "A")
	test("marzip", "lesser than", "marzipan")
	test("bobby", "not equal to", "suzzie")
	test("tyrone", "equal to", "tyrone")
}
