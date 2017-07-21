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
	testify "github.com/stretchr/testify/assert"
	"testing"
)

// FIX: things limiting reuse across tests:
// . concrete object pointers;
//   alt: get by name and return a nothing object
// . many "builders" use map, you cant re-make that map just by calling Rtm.New;
//   alt: revisit builders
// . object needs Objects for get/pointer, which ties objects to classes early;
//   alt: get rid of pointers/object lookup.
func TestGroupDefaults(t *testing.T) {
	assert := testify.New(t)

	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		objectList...)
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
		namePatterns,
	)
	assert.NoError(e)

	var lines printer.Span
	run := rtm.New(classes).Objects(objects).Patterns(patterns).Writer(&lines).Rtm()

	os := &core.Objects{nameList}
	// test the underlying grouping alg:
	if grps, e := MakeGroups(run, os); assert.NoError(e) {
		assert.Empty(grps.Grouped)
		assert.Len(grps.Ungrouped, 5)
	}
	// then test the actual output:
	prn := &PrintNondescriptObjects{os}
	if e := prn.Execute(run); assert.NoError(e) {
		assert.Equal("an empire apple, Mildred, a pen, a thing, and a thing", lines.String())
	}
}

func TestGroupPatterns(t *testing.T) {
	assert := testify.New(t)

	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		objectList...)
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
		namePatterns,
		groupPatterns,
	)
	assert.NoError(e)

	var lines printer.Span
	run := rtm.New(classes).Objects(objects).Patterns(patterns).Writer(&lines).Rtm()

	os := &core.Objects{nameList}
	// test the underlying grouping alg:
	if grps, e := MakeGroups(run, os); assert.NoError(e) {
		assert.Len(grps.Grouped, 1)

		for g, objs := range grps.Grouped {
			assert.False(g.Innumerable)
			assert.Equal("thing", g.Label)
			assert.Equal(GroupWithoutObjects, g.ObjectGrouping)
			assert.Len(objs.Objects, 2)
		}
		assert.Len(grps.Ungrouped, 3)
	}
	// then test the actual output:
	prn := &PrintNondescriptObjects{os}
	if e := prn.Execute(run); assert.NoError(e) {
		assert.Equal("an empire apple, Mildred, a pen, and two things", lines.String())
	}
}
