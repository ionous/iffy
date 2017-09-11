package rule_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func Int(i int) *core.Num {
	return &core.Num{float64(i)}
}

// Factorial computes an integer multiplied by the factorial of the integer below it.
type Factorial struct {
	Num float64
}

func TestFactorial(t *testing.T) {
	assert := testify.New(t)
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	unique.PanicBlocks(cmds,
		(*rule.Commands)(nil),
		(*core.Commands)(nil))
	patterns := unique.NewStack(classes)
	unique.PanicTypes(patterns,
		(*Factorial)(nil))
	assert.Contains(classes, ident.IdOf("Factorial"), "adding to patterns should add to classes")

	var root struct{ rule.Mandates }
	if c, ok := cmds.NewXBuilder(&root, core.Xform{}); ok {
		if c.Cmds().Begin() {
			if c.Cmd("number rule", "factorial").Begin() {
				c.Param("if").Cmd("compare num", c.Cmd("get", "@", "num"), c.Cmd("equal to"), 0)
				c.Param("decide").Val(1)
				c.End()
			}
			if c.Cmd("number rule", "factorial").Begin() {
				if c.Param("decide").Cmd("mul", c.Cmd("get", "@", "num")).Begin() {
					if c.Cmd("determine").Begin() {
						// alt: register factorial as a shadow class, and trigger a new factorial.
						// this currently relies on "set num" returning its "this":
						// therefore "determine" receives the factorial object, and re-runs the pattern.
						c.Cmd("set num", "@", "num", c.Cmd("sub", c.Cmd("get", "@", "num"), 1))
						c.End()
					}
					c.End()
				}
				c.End()
			}
			c.End()
		}
		rules := rule.MakeRules()

		if e := c.Build(); assert.NoError(e) {
			if els := root.Mandates; assert.Len(els, 2) {
				// test.Log(pretty.Sprint(els))
				if e := els.Mandate(patterns.Types, rules); e != nil {
					t.Fatal(e)
				}
			}
		}
		//
		objects := obj.NewObjects()
		run := rtm.New(classes).Objects(objects).Rules(rules).Rtm()
		peal := run.GetPatterns()
		if numberPatterns := peal.Numbers; assert.Len(numberPatterns, 1) {
			if factPattern := numberPatterns[ident.IdOf("factorial")]; assert.Len(factPattern, 2) {
				//
				if n, e := run.GetNumMatching(run, run.Emplace(&Factorial{3})); assert.NoError(e) {
					fac := 3 * (2 * (1 * 1))
					assert.EqualValues(fac, n)
				}
			}
		}
	}
}
