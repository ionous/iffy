package std

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std/group"
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
	// note: even without grouping rules, our article printing still gathers unnamed items.
	t.Run("no grouping", func(t *testing.T) {
		groupTest(t, "Mildred, an empire apple, a pen, and two other things",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2"),
			printPatterns)
	})
	//
	t.Run("default grouping", func(t *testing.T) {
		groupTest(t, "Mildred, an empire apple, a pen, and two things",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2"),
			printPatterns, group.GroupPatterns)
	})
	//
	several := func(c *ops.Builder) {
		if c.Cmd("run rule", "group together").Begin() {
			c.Param("if").Cmd("is same class", c.Cmd("get", "@", "target"), "thing")
			if c.Param("decide").Cmds().Begin() {
				c.Cmd("set bool", "@", "with articles", true)
				c.End()
			}
			c.End()
		}
		if c.Cmd("run rule", "print several").Begin() {
			c.Param("if").Cmd("all true", c.Cmds(
				c.Cmd("is same class", c.Cmd("get", "@", "target"), "thing"),
				c.Cmd("compare num", c.Cmd("get", "@", "group size"), c.Cmd("greater than"), 1)),
			)
			//
			if c.Param("decide").Cmds().Begin() {
				c.Cmd("print text", "a few things")
				c.End()
			}
			c.End() // print several
		}
	}
	t.Run("several", func(t *testing.T) {
		// note: we are grouping together the things, they become a single unit --
		// therefore no comma between Mildred the person and the single unit of all things.
		groupTest(t, "Mildred and an empire apple, a pen, and a few things",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2"),
			printPatterns, group.GroupPatterns, several)
	})
	t.Run("not several", func(t *testing.T) {
		groupTest(t, "Mildred and an empire apple, a pen, and one other thing",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1"),
			printPatterns, group.GroupPatterns, several)
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
			printPatterns, group.GroupPatterns, replacement)
	})
	// verify: if there arent multiple things to group, then we shouldnt be grouping
	t.Run("replacement none", func(t *testing.T) {
		groupTest(t, "Mildred and an empire apple",
			//
			sliceOf.String("mildred", "apple"),
			printPatterns, group.GroupPatterns, replacement)
	})
	// Before listing contents: group Scrabble pieces together.
	// Before grouping together Scrabble pieces, say "the tiles ".
	// After grouping together Scrabble pieces, say " from a Scrabble set".
	t.Run("fancy", func(t *testing.T) {
		fancy := func(c *ops.Builder) {
			if c.Cmd("run rule", "group together").Begin() {
				c.Param("if").Cmd("is same class", c.Cmd("get", "@", "target"), "scrabble tile")
				if c.Param("decide").Cmds().Begin() {
					c.Cmd("set text", "@", "label", "the tiles")
					c.Cmd("set bool", "@", "innumerable", true)
					c.Cmd("set bool", "@", "without articles", true)
					c.End()
				}
				c.End()
			}
			if c.Cmd("run rule", "print group").Begin() {
				c.Param("if").Cmd("compare text", c.Cmd("get", "@", "label"), c.Cmd("equal to"), "the tiles")
				c.Param("continue").Cmd("continue before")
				if c.Param("decide").Cmds().Begin() {
					c.Cmd("print text", "from a Scrabble set")
					c.End()
				}
				c.End()
			}
		}
		t.Run("tiles", func(t *testing.T) {
			groupTest(t, "the tiles X, W, F, Y, and Z from a Scrabble set",
				//
				sliceOf.String("x", "w", "f", "y", "z"),
				printPatterns, group.GroupPatterns, fancy)
		})
		t.Run("more", func(t *testing.T) {
			groupTest(t, "Mildred and the tiles X, W, F, Y, and Z from a Scrabble set",
				//
				sliceOf.String("mildred", "x", "w", "f", "y", "z"),
				printPatterns, group.GroupPatterns, fancy)
		})
	})
	// //group utensils together as "text".
	t.Run("parenthetical", func(t *testing.T) {
		groupTest(t, "Mildred and five scrabble tiles ( X, W, F, Y, and Z )",
			sliceOf.String("mildred", "x", "w", "f", "y", "z"),
			//
			printPatterns, group.GroupPatterns, func(c *ops.Builder) {
				if c.Cmd("run rule", "group together").Begin() {
					c.Param("if").Cmd("is same class", c.Cmd("get", "@", "target"), "scrabble tile")
					if c.Param("decide").Cmds().Begin() {
						c.Cmd("set text", "@", "label", "scrabble tiles")
						c.Cmd("set bool", "@", "without articles", true)
						c.End()
					}
					c.End()
				}
			})
	})
	t.Run("article parenthetical", func(t *testing.T) {
		groupTest(t, "Mildred and five scrabble tiles ( a X, a W, a F, a Y, and a Z )",
			//
			sliceOf.String("mildred", "x", "w", "f", "y", "z"),
			printPatterns, group.GroupPatterns, func(c *ops.Builder) {
				if c.Cmd("run rule", "group together").Begin() {
					c.Param("if").Cmd("is same class", c.Cmd("get", "@", "target"), "scrabble tile")
					if c.Param("decide").Cmds().Begin() {
						c.Cmd("set text", "@", "label", "scrabble tiles")
						c.Cmd("set bool", "@", "with articles", true)
						c.End()
					}
					c.End()
				}
			})
	})
	unnamedThings := func(c *ops.Builder) {
		if c.Cmd("run rule", "group together").Begin() {
			c.Param("if").Cmd("is same class", c.Cmd("get", "@", "target"), "thing")
			if c.Param("decide").Cmds().Begin() {
				c.Cmd("set text", "@", "label", "things")
				c.Cmd("set bool", "@", "with articles", true)
				c.End()
			}
			c.End()
		}
	}
	t.Run("edge one", func(t *testing.T) {
		groupTest(t, "Mildred and three things ( an empire apple, a pen, and one other thing )",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1"),
			printPatterns, group.GroupPatterns, unnamedThings)
	})
	t.Run("edge two", func(t *testing.T) {
		groupTest(t, "Mildred and four things ( an empire apple, a pen, and two other things )",
			//
			sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2"),
			printPatterns, group.GroupPatterns, unnamedThings)
	})
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
	if assert.NoError(e) {
		var lines printer.Span
		run := rtm.New(classes).Objects(objects).Patterns(patterns).Writer(&lines).Rtm()

		prn := &PrintNondescriptObjects{&core.Objects{names}}
		if e := prn.Execute(run); assert.NoError(e) {
			if assert.Equal(match, lines.String()) {
				t.Log("matched:", match)
			}
		}
	}
}
