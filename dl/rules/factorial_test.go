package rules_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

// Factorial computes an integer multiplied by the factorial of the integer below it.
type Factorial struct {
	Num float64
}

func TestFactorial(t *testing.T) {
	assert := testify.New(t)
	//
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	patterns := unique.NewStack(cmds.ShadowTypes)
	contract := pat.MakeContract(patterns.Types)

	unique.PanicBlocks(cmds,
		(*rules.Commands)(nil),
		(*core.Commands)(nil))
	unique.PanicTypes(patterns,
		(*Factorial)(nil))
	assert.Contains(classes, ident.IdOf("Factorial"), "adding to patterns should add to classes")

	var root struct{ rules.Mandates }
	c := cmds.NewBuilder(&root, core.Xform{})
	if c.Cmd("number rule", "factorial").Begin() {
		c.Param("if").Cmd("compare num", c.Cmd("get", "@", "num"), c.Cmd("equal to"), 0)
		c.Param("decide").Val(1)
		c.End()
	}
	if c.Cmd("number rule", "factorial").Begin() {
		if c.Param("decide").Cmd("mul", c.Cmd("get", "@", "num")).Begin() {
			c.Cmd("determine", c.Cmd("factorial", c.Cmd("sub", c.Cmd("get", "@", "num"), 1)))
			c.End()
		}
		c.End()
	}
	//
	if e := c.Build(); assert.NoError(e) {
		if els := root.Mandates; assert.Len(els, 2) {
			// test.Log(pretty.Sprint(els))
			if e := els.Mandate(contract); e != nil {
				t.Fatal(e)
			}
		}
	}
	if run, e := rtm.New(classes).Rules(contract).Rtm(); assert.NoError(e) {
		if n, e := run.GetNumMatching(run.Emplace(&Factorial{3})); assert.NoError(e) {
			fac := 3 * (2 * (1 * 1))
			assert.EqualValues(fac, n)
		}
	}
}
