package core_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ops"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"github.com/stretchr/testify/suite"
	r "reflect"
	"strings"
	"testing"
)

func TestCore(t *testing.T) {
	suite.Run(t, new(CoreSuite))
}

// Regular expression to select test suites specified command-line argument "-run". Regular expression to select the methods of test suites specified command-line argument "-m"
type CoreSuite struct {
	suite.Suite
	ops *ops.Ops

	run   rt.Runtime
	lines rtm.LineWriter

	classes   *ref.Classes
	objects   *ref.Objects
	relations *ref.Relations

	unique *unique.Objects
}

func (assert *CoreSuite) Lines() (ret []string) {
	ret = assert.lines.Lines()
	assert.lines = rtm.LineWriter{}
	return
}

func (assert *CoreSuite) SetupTest() {
	errutil.Panic = false
	assert.ops = ops.NewOps((*core.Commands)(nil))

	assert.unique = unique.NewObjects()
	assert.classes = ref.NewClasses()
	assert.objects = ref.NewObjects(assert.classes)
	assert.relations = ref.NewRelations(assert.classes, assert.objects)

	unique.RegisterTypes(unique.PanicTypes(assert.unique),
		(*core.CycleCounter)(nil),
		(*core.ShuffleCounter)(nil),
		(*core.StoppingCounter)(nil))
}

func (assert *CoreSuite) newRuntime(c *ops.Builder) (ret rt.Runtime, err error) {
	if _, e := c.Build(); e != nil {
		err = e
	} else {

		mm := unique.PanicTypes(assert.classes)
		unique.RegisterTypes(mm,
			(*core.NumberCounter)(nil),
			(*core.TextCounter)(nil))

		// add all the helper classes we registered via unique
		for _, rtype := range assert.unique.Types {
			unique.RegisterTypes(mm, r.New(rtype).Interface())
		}

		if inst, e := assert.unique.Generate(); e != nil {
			err = e
		} else {
			unique.RegisterValues(unique.PanicValues(assert.objects), inst...)

			run := rtm.NewRtm(assert.classes, assert.objects, assert.relations)
			run.PushWriter(&assert.lines)
			ret = run
		}
	}
	return
}

func (assert *CoreSuite) matchFunc(
	build func(c *ops.Builder),
	compare func(expected []string),
) {
	var root struct{ Eval rt.Execute }
	if c, ok := assert.ops.NewBuilder(&root); ok {
		build(c)
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if e := root.Eval.Execute(run); assert.NoError(e) {
				compare(assert.Lines())
			}
		}
	}
}

func (assert *CoreSuite) match(
	build func(c *ops.Builder), expected ...string) {
	assert.matchFunc(build, func(lines []string) {
		assert.Equal(expected, lines)
	})
}

func (assert *CoreSuite) TestShortcuts() {
	var root struct {
		Eval rt.TextEval
	}
	if c, ok := assert.ops.NewBuilder(&root); ok {
		c.Val("shortcut")
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if res, e := root.Eval.GetText(run); assert.NoError(e) {
				assert.EqualValues("shortcut", res)
			}
		}
	}
}

// TestAllTrue ensure AllTrue operates on boolean literals as "and".
func (assert *CoreSuite) TestAllTrue() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a, b, res bool) {
		if c, ok := assert.ops.NewBuilder(&root); ok {
			c.Cmd("all true", c.Cmds(
				c.Cmd("bool", a),
				c.Cmd("bool", b)))
			//
			if run, e := assert.newRuntime(c); assert.NoError(e) {
				if ok, e := root.Eval.GetBool(run); assert.NoError(e) {
					assert.EqualValues(res, ok)
				}
			}
		}
	}
	// ******
	test(true, false, false)
	test(true, true, true)
	test(false, false, false)
}

// TestAnyTrue ensure AnyTrue operates on boolean literals as "or".
func (assert *CoreSuite) TestAnyTrue() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a, b, res bool) {
		if c, ok := assert.ops.NewBuilder(&root); ok {
			if c.Cmd("any true").Begin() {
				c.Cmds(c.Cmd("bool", a), c.Cmd("bool", b))
				c.End()
			}
			// /
			if run, e := assert.newRuntime(c); assert.NoError(e) {
				if ok, e := root.Eval.GetBool(run); assert.NoError(e) {
					assert.EqualValues(res, ok)
				}
			}
		}
	}
	test(true, false, true)
	test(true, true, true)
	test(false, false, false)
}

func (assert *CoreSuite) TestCompareNum() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a float64, op string, b float64) {
		if c, ok := assert.ops.NewBuilder(&root); ok {
			if c.Cmd("compare num").Begin() {
				c.Val(a).Cmd(op).Val(b)
				c.End()
			}

			if run, e := assert.newRuntime(c); assert.NoError(e) {
				if ok, e := root.Eval.GetBool(run); assert.NoError(e) {
					assert.True(ok)
				}
			}
		}
	}
	test(10, "greater than", 1)
	test(1, "lesser than", 10)
	test(10, "not equal to", 1)
	test(10, "equal to", 10)
}

func (assert *CoreSuite) TestCompareText() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a, op, b string) {
		if c, ok := assert.ops.NewBuilder(&root); ok {
			c.Cmd("compare text", c.Val(a), c.Cmd(op), c.Val(b))
			//
			if run, e := assert.newRuntime(c); assert.NoError(e) {
				if ok, e := root.Eval.GetBool(run); assert.NoError(e) {
					assert.True(ok, strings.Join([]string{a, op, b}, " "))
				}
			}
		}
	}
	test("Z", "greater than", "A")
	test("marzip", "lesser than", "marzipan")
	test("bobby", "not equal to", "suzzie")
	test("tyrone", "equal to", "tyrone")
}
