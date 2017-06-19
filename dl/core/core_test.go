package core

import (
	. "github.com/ionous/iffy/dl/tests"
	"github.com/ionous/iffy/ops"
	"github.com/ionous/iffy/reflector"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

var people = []interface{}{
	&Person{"alice", Listening, Standing},
	&Person{"bob", Listening, Sitting},
	&Person{"carol", Laughing, Sitting},
}

type CoreSuite struct {
	suite.Suite
	ops *ops.Ops
	run rt.Runtime
}

func TestCoreSuite(t *testing.T) {
	suite.Run(t, new(CoreSuite))
}

func (t *CoreSuite) Log(args ...interface{}) {
	t.T().Log(args...)
}

func (t *CoreSuite) SetupTest() {
	t.ops = ops.NewOps((*Commands)(nil))
	if m, e := reflector.MakeModel(); e != nil {
		panic(e)
	} else {
		t.run = rtm.NewRtm(m)
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

// if m, e := reflector.MakeModel(people...); t.NoError(e) {
// 	rtm.NewRtm(m)
// 	{
// 		var root struct {
// 			Eval rt.BoolEval
// 		}
// 		if c := ops.Build(&root); c.Args {
// 			if c := c.Cmd("all true").Array(); c.Cmds {
// 				c.Cmd("get", "alice", "listening")
// 				c.Cmd("get", "bob", "sitting")
// 				c.Cmd("get", "carol", "laughing")
// 			}
// 		}
// 		// for _, op := range ops {
// 		// }
// 	}
// }

// FIX: test all forms of Get
// TODO: test literals as eval subsitution
