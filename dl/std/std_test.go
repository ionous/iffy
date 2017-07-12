package std_test

import (
	"github.com/ionous/iffy/dl/core"
	. "github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/pat/patspec"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	// "github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"testing"
)

// go test -run ''      # Run all tests.
// go test -run Foo     # Run top-level tests matching "Foo", such as "TestFooBar".
// go test -run Foo/A=  # For top-level tests matching "Foo", run subtests matching "A=".
// go test -run /A=1    # For all top-level tests, run subtests matching "A=1".
func TestStd(t *testing.T) {
	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))

	objects := ref.NewObjects(classes)
	// an unnamed thing
	unnamedThing := &Thing{}
	// a named thing
	genericPen := &Thing{
		Kind: Kind{Name: "pen"},
	}
	// a thing with a printed name
	plasticSword := &Thing{
		Kind: Kind{Name: "sword", PrintedName: "plastic sword"},
	}
	unique.RegisterValues(unique.PanicValues(objects),
		unnamedThing, genericPen, plasticSword)

	patterns := patbuilder.NewPatterns(classes)
	unique.RegisterBlocks(unique.PanicTypes(patterns),
		(*Patterns)(nil))

	ops := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(ops),
		(*core.Commands)(nil),
		(*patspec.Commands)(nil),
	)

	t.Run("Names", func(t *testing.T) {
		assert := assert.New(t)
		var root struct {
			Patterns patspec.PatternSpecs
		}
		if c, ok := ops.NewBuilder(&root); assert.True(ok) {
			if c.Param("patterns").Cmds().Begin() {
				buildPatterns(c)
				c.End()
			}
			_, e := c.Build()
			assert.NoError(e)
		}
		// t.Log("patterns", pretty.Sprint(root))
		if e := root.Patterns.Generate(patterns); e != nil {
			t.Fatal(e)
		}
		// FIX: so, we still cant have runtime ops to test these patterns,
		// because we cant zip things
		run := rtm.New(classes).Objects(objects).Patterns(patterns).NewRtm()
		//
		match(run, assert, "plastic sword", &plasticSword.Kind)
		match(run, assert, "pen", &genericPen.Kind)
		match(run, assert, "thing", &unnamedThing.Kind)
	})

	// <tear-down code>

}

func match(run rt.Runtime, assert *assert.Assertions, match string, kind *Kind) (okay bool) {
	var lines rtm.LineWriter
	run.PushWriter(&lines)
	defer run.PopWriter()
	if printName, e := run.Emplace(&PrintName{kind}); assert.NoError(e) {
		if _, e := run.ExecuteMatching(printName); assert.NoError(e) {
			okay = assert.EqualValues(sliceOf.String(match), lines.Lines())
		}
	}
	return
}

func buildPatterns(c *ops.Builder) {
	// its a little heavy to do this with patterns, but -- its a good test of the system.
	// print the class name if all else fails
	if c.Cmd("run rule", "print name").Begin() {
		c.Param("decide").Cmd("print text", c.Cmd("class name", c.Cmd("get", "@", "target")))
		c.End()
	}
	// prefer the object name, so long as it was specified by the user.
	if c.Cmd("run rule", "print name").Begin() {
		// # is used only for system names, not user author names.
		c.Param("if").Cmd("is not", c.Cmd("includes", c.Cmd("get", c.Cmd("get", "@", "target"), "name"), "#"))
		c.Param("decide").Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "name"))
		c.End()
	}
	// perfer the printed name above all else
	if c.Cmd("run rule", "print name").Begin() {
		c.Param("if").Cmd("is not", c.Cmd("is empty", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name")))
		c.Param("decide").Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name"))
		c.End()
	}

}
