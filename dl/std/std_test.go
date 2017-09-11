package std_test

import (
	"github.com/ionous/iffy/dl/core"
	. "github.com/ionous/iffy/dl/std"

	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestStd(t *testing.T) {
	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOpsX(classes, core.Xform{})    // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types

	unique.PanicBlocks(classes,
		(*Classes)(nil))

	unique.PanicBlocks(patterns,
		(*Patterns)(nil))

	objects := obj.NewObjects()
	unique.PanicValues(objects,
		Thingaverse.objects(sliceOf.String("apple", "pen", "thing#1", "thing#2"))...)

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*rule.Commands)(nil),
		(*Commands)(nil),
	)

	t.Run("Names", func(t *testing.T) {
		assert := testify.New(t)

		//
		rules, e := rule.Master(cmds, patterns, PrintNameRules)
		assert.NoError(e)

		// TODO: add test for: Rule for printing the name of the pen while taking inventory: say "useful pen".
		// TODO: add test for: A novel is a kind of thing. Dr Zhivago and Persuasion are novels. Before printing the name of a novel, say "[italic type]". After printing the name of a novel, say "[roman type]".”
		run := rtm.New(classes).Objects(objects).Rules(rules).Rtm()
		//
		t.Run("printed name", func(t *testing.T) {
			apple, _ := run.GetObject("apple")
			match(run, testify.New(t), "empire apple", &PrintName{apple.Id()})
		})
		t.Run("named", func(t *testing.T) {
			genericPen, _ := run.GetObject("pen")
			match(run, testify.New(t), "pen", &PrintName{genericPen.Id()})
		})
		t.Run("unnamed", func(t *testing.T) {
			unnamedThing, _ := run.GetObject("thing#1")
			match(run, testify.New(t), "thing", &PrintName{unnamedThing.Id()})
		})
		//
		t.Run("plural printed name", func(t *testing.T) {
			apple, _ := run.GetObject("apple")
			match(run, testify.New(t), "empire apples", &PrintPluralName{apple.Id()})
		})
		t.Run("plural named", func(t *testing.T) {
			genericPen, _ := run.GetObject("pen")
			match(run, testify.New(t), "pens", &PrintPluralName{genericPen.Id()})
		})
		t.Run("plural unnamed", func(t *testing.T) {
			unnamedThing, _ := run.GetObject("thing#1")
			match(run, testify.New(t), "things", &PrintPluralName{unnamedThing.Id()})
		})
		//
		t.Run("printed plural name", func(t *testing.T) {
			forced := "party favors"
			apple, _ := run.GetObject("apple")
			if e := apple.SetValue("printed plural name", forced); e != nil {
				t.Fatal(e)
			}

			match(run, testify.New(t), forced, &PrintPluralName{apple.Id()})
		})
		//
	})
}

func match(run rt.Runtime, assert *testify.Assertions, match string, op interface{}) (okay bool) {
	var lines printer.Lines
	run = rt.Writer(run, &lines)
	if e := run.ExecuteMatching(run, run.Emplace(op)); assert.NoError(e) {
		okay = assert.EqualValues(sliceOf.String(match), lines.Lines())
	}
	return
}
