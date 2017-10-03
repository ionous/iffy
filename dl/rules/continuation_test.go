package rules_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

// Order computes an integer multiplied by the factorial of the integer below it.
func TestOrder(t *testing.T) {
	assert := testify.New(t)
	//
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	patterns := unique.NewStack(cmds.ShadowTypes)
	contract := pat.MakeContract(patterns.Types)

	type Order struct{}

	unique.PanicBlocks(cmds,
		(*rules.Commands)(nil),
	// (*core.Commands)(nil))
	)
	unique.PanicTypes(patterns,
		(*Order)(nil))
	assert.Contains(classes, ident.IdOf("Order"), "adding to patterns should add to classes")

	var root struct{ rules.Mandates }
	c := cmds.NewBuilder(&root, core.Xform{})
	if c.Cmds().Begin() {
		if c.Cmd("list text", "order").Begin() {
			c.Param("decide").Val("a")
			c.End()
		}
		if c.Cmd("list text", "order").Begin() {
			c.Param("decide").Val("b")
			c.Param("continue").Cmd("continue after")
			c.End()
		}
		if c.Cmd("list text", "order").Begin() {
			c.Param("decide").Val("c")
			c.Param("continue").Cmd("continue before")
			c.End()
		}
		c.End()
	}
	//
	if e := c.Build(); assert.NoError(e) {
		if e := root.Mandate(contract); e != nil {
			t.Fatal(e)
		}
	}
	//
	if run, e := rtm.New(classes).Rules(contract).Rtm(); e != nil {
		t.Fatal("couldnt create runtime", e)
	} else {
		t.Log(pretty.Sprint(run.Rules))
		if n, e := run.GetTextStreamMatching(run.Emplace(&Order{})); e != nil {
			t.Fatal("no rules", e)
		} else {
			var res string
			expect := "bac"
			for n.HasNext() {
				if x, e := n.GetText(); e != nil {
					t.Fatal("error getting next", e)
				} else if len(x) == 0 {
					t.Fatal("empty result")
				} else if len(res) > len(expect) {
					t.Fatal("didnt end soon enough", res, x)
				} else {
					res += x
				}
			}
			if res != expect {
				t.Fatal("mismatched", res)
			}
		}
	}
}
