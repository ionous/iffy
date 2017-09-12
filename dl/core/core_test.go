package core_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
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
	cmds *ops.Ops

	run   rt.Runtime
	lines printer.Lines

	classes unique.Types
	objects *obj.ObjBuilder

	gen *unique.Objects
}

func (assert *CoreSuite) Lines() (ret []string) {
	ret = assert.lines.Lines()
	assert.lines = printer.Lines{}
	return
}

func (assert *CoreSuite) SetupTest() {
	errutil.Panic = false
	classes := make(unique.Types)
	assert.cmds = ops.NewOps(classes)
	unique.PanicBlocks(assert.cmds,
		(*core.Commands)(nil))

	assert.gen = unique.NewObjectGenerator()
	assert.classes = classes
	assert.objects = obj.NewObjects()

	unique.PanicBlocks(assert.gen,
		(*core.Counters)(nil))

	unique.PanicBlocks(assert.classes,
		(*core.Classes)(nil),
		(*core.Counters)(nil),
	)
}

func (assert *CoreSuite) newRuntime(c *ops.Builder) (ret rt.Runtime, err error) {
	if e := c.Build(); e != nil {
		err = e
	} else {
		if objs, e := assert.gen.Generate(); e != nil {
			err = e
		} else {
			unique.PanicValues(assert.objects, objs...)
			ret = rtm.New(assert.classes).Objects(assert.objects).Writer(&assert.lines).Rtm()
		}
	}
	return
}

func (assert *CoreSuite) matchFunc(
	build func(c spec.Block),
	compare func(expected []string),
) {
	var root struct{ Eval rt.Execute }
	c := assert.cmds.NewBuilder(&root, core.Xform{})
	build(c)
	if run, e := assert.newRuntime(c); assert.NoError(e) {
		if e := root.Eval.Execute(run); assert.NoError(e) {
			compare(assert.Lines())
		}
	}
}

func (assert *CoreSuite) matchLine(expected string,
	build func(c spec.Block)) {
	assert.matchLines([]string{expected}, build)
}

func (assert *CoreSuite) matchLines(expected []string,
	build func(c spec.Block)) {
	assert.matchFunc(build, func(lines []string) {
		assert.Equal(expected, lines)
	})
}

func (assert *CoreSuite) TestShortcuts() {
	var root struct {
		Eval rt.TextEval
	}
	c := assert.cmds.NewBuilder(&root, core.Xform{})
	c.Val("shortcut")
	if run, e := assert.newRuntime(c); assert.NoError(e) {
		if res, e := root.Eval.GetText(run); assert.NoError(e) {
			assert.EqualValues("shortcut", res)
		}
	}

}

// TestAllTrue ensure AllTrue operates on boolean literals as "and".
func (assert *CoreSuite) TestAllTrue() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a, b, res bool) {
		c := assert.cmds.NewBuilder(&root, core.Xform{})
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
		c := assert.cmds.NewBuilder(&root, core.Xform{})
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
	test(true, false, true)
	test(true, true, true)
	test(false, false, false)
}

func (assert *CoreSuite) TestCompareNum() {
	var root struct {
		Eval rt.BoolEval
	}
	test := func(a float64, op string, b float64) {
		c := assert.cmds.NewBuilder(&root, core.Xform{})
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
		c := assert.cmds.NewBuilder(&root, core.Xform{})
		c.Cmd("compare text", c.Val(a), c.Cmd(op), c.Val(b))
		//
		if run, e := assert.newRuntime(c); assert.NoError(e) {
			if ok, e := root.Eval.GetBool(run); assert.NoError(e) {
				assert.True(ok, strings.Join([]string{a, op, b}, " "))
			}
		}

	}
	test("Z", "greater than", "A")
	test("marzip", "lesser than", "marzipan")
	test("bobby", "not equal to", "suzzie")
	test("tyrone", "equal to", "tyrone")
}
