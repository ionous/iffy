package patspec_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/pat/patspec"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	// "github.com/kr/pretty"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPattern(assert *testing.T) {
	suite.Run(assert, new(PatternSuite))
}

type PatternSuite struct {
	suite.Suite
	ops *ops.Ops
}

func (assert *PatternSuite) SetupTest() {
	ops := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(ops),
		(*patspec.Commands)(nil),
		(*core.Commands)(nil))
	assert.ops = ops
}

func Int(i int) *core.Num {
	return &core.Num{float64(i)}
}

// Factorial computes an integer multiplied by the factorial of the integer below it.
type Factorial struct {
	Num float64
}

func (assert *PatternSuite) TestFactorial() {
	classes := ref.NewClasses()
	patterns := patbuilder.NewPatterns(classes)
	unique.RegisterTypes(unique.PanicTypes(patterns),
		(*Factorial)(nil))
	assert.Contains(classes.ClassMap, id.MakeId("Factorial"), "adding to patterns should add to classes")

	var root struct {
		Els patspec.PatternSpecs
	}
	if c, ok := assert.ops.NewBuilder(&root); ok {
		if c.Cmds().Begin() {
			if c.Cmd("number rule", "factorial").Begin() {
				// FIX? re: "equal to" - can literally detect string and make empty command?
				c.Param("if").Cmd("compare num", c.Cmd("get", "@", "num"), c.Cmd("equal to"), 0)
				c.Param("decide").Val(1)
				c.End()
			}
			if c.Cmd("number rule", "factorial").Begin() {
				if c.Param("decide").Cmd("mul", c.Cmd("get", "@", "num")).Begin() {
					if c.Cmd("determine").Begin() {
						// FIX: we need to be able to construct a factorial object from scratch
						// treating it just like it were any other command
						c.Cmd("set num", "@", "num", c.Cmd("sub", c.Cmd("get", "@", "num"), 1))
						c.End()
					}
					c.End()
				}
				c.End()
			}
			c.End()
		}
		//
		test := assert.T()
		if e := c.Build(); assert.NoError(e) {
			if els := root.Els; assert.Len(els, 2) {
				// test.Log(pretty.Sprint(els))
				if e := els.Generate(patterns); e != nil {
					test.Fatal(e)
				}
			}
		}
		//
		objects := ref.NewObjects(classes)
		run := rtm.New(classes).Objects(objects).Patterns(patterns).Rtm()
		peal := run.GetPatterns()
		if numberPatterns := peal.NumberMap; assert.Len(numberPatterns, 1) {
			if factPattern := numberPatterns[id.MakeId("factorial")]; assert.Len(factPattern, 2) {
				//
				if fact, e := run.Emplace(&Factorial{3}); assert.NoError(e) {
					if n, e := run.GetNumMatching(fact); assert.NoError(e) {
						fac := 3 * (2 * (1 * 1))
						assert.EqualValues(fac, n)
					}
				}
			}
		}
	}
}
