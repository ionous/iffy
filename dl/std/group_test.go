package std

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/pat/patspec"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestGrouping(t *testing.T) {
	t.Run("no grouping", func(t *testing.T) {
		groupTest(t, "Mildred, an empire apple, a pen, a thing, and a thing",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2"),
			namePatterns)
	})
	t.Run("default grouping", func(t *testing.T) {
		groupTest(t, "Mildred, an empire apple, a pen, and two things",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2"),
			namePatterns, groupPatterns)
	})
	// Rule for grouping together utensils: say "the usual utensils".
	replacement := func(c *ops.Builder) {
		if c.Cmd("run rule", "group together").Begin() {
			c.Param("if").Cmd("is same class", c.Cmd("get", "@", "target"), "thing")
			if c.Param("decide").Cmds().Begin() {
				c.Cmd("set text", "@", "label", "some things")
				c.Cmd("set bool", "@", "innumerable", true)
				c.Cmd("set bool", "@", "without objects", true)
				c.End()
			}
			c.End()
		}
	}
	t.Run("replacement text", func(t *testing.T) {
		groupTest(t, "Mildred and some things",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2"),
			namePatterns, groupPatterns, replacement)
	})
	// if there arent multiple things to group, then we shouldnt be grouping
	t.Run("replacement none", func(t *testing.T) {
		groupTest(t, "Mildred and an empire apple",
			//
			sliceOf.String("mildred", "apple"),
			namePatterns, groupPatterns, replacement)
	})
	// t.Run("fancy grouping", func(t *testing.T) {
	// 	groupTest(t,
	// 		"Mildred, and the tiles X, W, F, Y and Z from a Scrabble set",
	// 		sliceOf.String("mildred", "x", "w", "f", "y", "z"),
	// 		namePatterns, groupPatterns)
	// })
	// //group utensils together as "text".
	// t.Run("parenthetical", func(t *testing.T) {
	// 	groupTest(t, "Mildred, and five scrabble tiles ( X, W, F, Y and Z )",
	// 		sliceOf.String("mildred", "x", "w", "f", "y", "z"),

	// 		namePatterns, groupPatterns)
	// })
	// t.Run("article parenthetical", func(t *testing.T) {
	// 	groupTest(t, "Mildred, and five scrabble tiles ( a X, a W, a F, a Y and a Z )",
	// 		sliceOf.String("mildred", "x", "w", "f", "y", "z"),
	// 		namePatterns, groupPatterns)
	// })
	// t.Run("edge case", func(t *testing.T) {
	// 	groupTest(t, "Mildred, and four things ( a pen, an empire apple, and two things )",
	// 		sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2"),
	// 		namePatterns, groupPatterns)
	// })
}

// FIX: things limiting reuse across tests:
// . concrete object pointers;
//   alt: get by name and return a nothing object
// . many "builders" use map, you cant re-make that map just by calling Rtm.New;
//   alt: revisit builders
// . object needs Objects for get/pointer, which ties objects to classes early;
//   alt: get rid of pointers/object lookup.
func groupTest(t *testing.T, match string, names []string, patternSpec ...func(*ops.Builder)) {
	assert := testify.New(t)

	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))
	// for testing:
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*ScrabbleTile)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		Thingaverse.objects(names)...)

	cmds := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*core.Commands)(nil),
		(*Commands)(nil),
		(*patspec.Commands)(nil),
	)
	unique.RegisterBlocks(unique.PanicTypes(cmds.ShadowTypes),
		(*Patterns)(nil),
	)
	//
	patterns, e := patbuilder.NewPatternMaster(cmds, classes,
		(*Patterns)(nil)).Build(
		patternSpec...,
	)
	assert.NoError(e)

	var lines printer.Span
	run := rtm.New(classes).Objects(objects).Patterns(patterns).Writer(&lines).Rtm()

	prn := &PrintNondescriptObjects{&core.Objects{names}}
	if e := prn.Execute(run); assert.NoError(e) {
		assert.Equal(match, lines.String())
	}
}
