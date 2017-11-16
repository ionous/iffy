package std_test

import (
	"bytes"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/express"
	"github.com/ionous/iffy/dl/rules"
	. "github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestStd(t *testing.T) {
	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOps(classes)                   // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types

	unique.PanicBlocks(classes,
		(*Classes)(nil))

	unique.PanicBlocks(patterns,
		(*Patterns)(nil))

	var objects obj.Registry
	objects.RegisterValues(Thingaverse.objects(sliceOf.String(
		"apple", "pen", "thing#1", "thing#2")))

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*rules.Commands)(nil),
		(*Commands)(nil),
	)

	t.Run("Names", func(t *testing.T) {
		assert := testify.New(t)
		//
		xform := express.NewTransform(cmds, nil)
		rules, e := rules.Master(cmds, xform, patterns, PrintNameRules)
		assert.NoError(e)

		// TODO: add test for: Rule for printing the name of the pen while taking inventory: say "useful pen".
		// TODO: add test for: A novel is a kind of thing. Dr Zhivago and Persuasion are novels. Before printing the name of a novel, say "[italic type]". After printing the name of a novel, say "[roman type]".”
		run, e := rtm.New(classes).Objects(objects).Rules(rules).Rtm()
		assert.NoError(e)
		//
		match := func(t *testing.T, what string, op interface{}) {
			var buf bytes.Buffer
			if e := rt.WritersBlock(run, &buf, func() error {
				return run.ExecuteMatching(run.Emplace(op))
			}); e != nil {
				t.Fatal(e)
			} else if res := buf.String(); res != what {
				t.Fatalf("%s != %s", res, what)
			}
		}

		t.Run("printed name", func(t *testing.T) {
			apple, _ := run.GetObject("apple")
			match(t, "empire apple", &PrintName{apple.Id()})
		})
		t.Run("named", func(t *testing.T) {
			genericPen, _ := run.GetObject("pen")
			match(t, "pen", &PrintName{genericPen.Id()})
		})
		t.Run("unnamed", func(t *testing.T) {
			unnamedThing, _ := run.GetObject("thing#1")
			match(t, "thing", &PrintName{unnamedThing.Id()})
		})
		//
		t.Run("plural printed name", func(t *testing.T) {
			apple, _ := run.GetObject("apple")
			match(t, "empire apples", &PrintPluralName{apple.Id()})
		})
		t.Run("plural named", func(t *testing.T) {
			genericPen, _ := run.GetObject("pen")
			match(t, "pens", &PrintPluralName{genericPen.Id()})
		})
		t.Run("plural unnamed", func(t *testing.T) {
			unnamedThing, _ := run.GetObject("thing#1")
			match(t, "things", &PrintPluralName{unnamedThing.Id()})
		})
		//
		t.Run("printed plural name", func(t *testing.T) {
			forced := "party favors"
			apple, _ := run.GetObject("apple")
			if e := apple.SetValue("printed plural name", forced); e != nil {
				t.Fatal(e)
			}

			match(t, forced, &PrintPluralName{apple.Id()})
		})
		//
	})
}
