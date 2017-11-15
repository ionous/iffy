package express_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/express"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

// TestOps to verify iffy statements.
// . text expressions: tested elsewhere
// . number, etc. expressions
// . direct calls: "{Story.score|testScore!}"
// . pattern determinations: c.Cmd("set text", "story", "status left", "{playerSurroundings!}")
// . pattern evaluations:  c.Cmd("{print banner text!}"
//
func TestOps(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	patterns := unique.NewStack(cmds.ShadowTypes)

	// a sample command.
	type TestScore struct{ rt.NumberEval }
	// a sample, empty, pattern.
	type TestPattern struct {
		Story ident.Id `if:"cls:kind"`
	}

	unique.PanicTypes(cmds,
		(*TestScore)(nil))
	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*rules.RuntimeCmds)(nil), // determine.
		(*express.Commands)(nil))
	unique.PanicTypes(patterns,
		(*TestPattern)(nil))

	xform := express.NewTransform(cmds, nil)

	t.Run("property", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, xform)
		//
		// HOW DOES THIS EVEN WORK!?!
		//
		c.Cmd("test score", "{5 + 5}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{
				NumberEval: &core.Add{
					A: &core.Num{Num: 5},
					B: &core.Num{Num: 5},
				},
			}, root.NumberEval)
		}
	})
	//
	t.Run("global", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("test score", "{Story.score}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			//
			//  IF THE OTHER BIT WORKS, WHY DO WE NEED RENDER?
			//  AT THE VERY LEAST WHY NOT GetAt
			//
			testEqual(t, &TestScore{&express.Render{
				&core.Object{"Story"}, "score"},
			}, root.NumberEval)
		}
	})
	//
	t.Run("run pipe", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("{Story.score|testScore!}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{&express.Render{
				&core.Object{"Story"}, "score"},
			}, root.NumberEval)
		}
	})
	t.Run("run params", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("{testScore! Story.score}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{&express.Render{
				&core.Object{"Story"}, "score"},
			}, root.NumberEval)
		}
	})
	// TESTPATTERN/1
	t.Run("determine", func(t *testing.T) {
		var root struct{ rt.TextEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("{testPattern! story}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testStrings(t, root.TextEval, `&rules.Determine{Obj:MakeTestPattern{Story:&express.GetAt{Name:"story"}}}`)
		}
	})
}

func testEqual(t *testing.T, expect, res interface{}) {
	t.Log("got:", pretty.Sprint(res))
	if !testify.ObjectsAreEqualValues(expect, res) {
		t.Log(pretty.Diff(res, expect))
		t.Log("want:", pretty.Sprint(expect))
		t.FailNow()
	}
}

func testStrings(t *testing.T, res interface{}, want string) {
	got := pretty.Sprintf("%#v", res)
	t.Log("got:", got)
	if got != want {
		t.Fatalf("want: %s", want)
	}
}
