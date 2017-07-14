package std_test

import (
	"github.com/ionous/iffy/dl/core"
	. "github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/dl/text"
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
		(*text.Commands)(nil),
		(*patspec.Commands)(nil),
	)

	t.Run("Names", func(t *testing.T) {
		_assert := assert.New(t)
		var root struct {
			Patterns patspec.PatternSpecs
		}
		if c, ok := ops.NewBuilder(&root); _assert.True(ok) {
			if c.Param("patterns").Cmds().Begin() {
				buildPatterns(c)
				c.End()
			}
			e := c.Build()
			t.Fatal(e)
		}
		// t.Log("patterns", pretty.Sprint(root))
		if e := root.Patterns.Generate(patterns); e != nil {
			t.Fatal(e)
		}
		// FIX: so, we still cant have runtime ops to test these patterns,
		// because we cant zip things
		run := rtm.New(classes).Objects(objects).Patterns(patterns).NewRtm()
		//
		t.Run("printed name", func(t *testing.T) {
			match(run, assert.New(t), "plastic sword", &PrintName{&plasticSword.Kind})
		})
		t.Run("named", func(t *testing.T) {
			match(run, assert.New(t), "pen", &PrintName{&genericPen.Kind})
		})
		t.Run("unnamed", func(t *testing.T) {
			match(run, assert.New(t), "thing", &PrintName{&unnamedThing.Kind})
		})

		//
		t.Run("plural printed name", func(t *testing.T) {
			match(run, assert.New(t), "plastic swords", &PrintPluralName{&plasticSword.Kind})
		})
		t.Run("plural named", func(t *testing.T) {
			match(run, assert.New(t), "pens", &PrintPluralName{&genericPen.Kind})
		})
		t.Run("plural unnamed", func(t *testing.T) {
			match(run, assert.New(t), "things", &PrintPluralName{&unnamedThing.Kind})
		})
		//
		forced := "party favors"
		plasticSword.PrintedPluralName = forced
		t.Run("printed plural name", func(t *testing.T) {
			match(run, assert.New(t), forced, &PrintPluralName{&plasticSword.Kind})
		})
	})
	// <tear-down code>

}

func match(run rt.Runtime, _assert *assert.Assertions, match string, op interface{}) (okay bool) {
	var lines rtm.LineWriter
	run.PushWriter(&lines)
	defer run.PopWriter()
	if printName, e := run.Emplace(op); _assert.NoError(e) {
		if _, e := run.ExecuteMatching(printName); _assert.NoError(e) {
			okay = _assert.EqualValues(sliceOf.String(match), lines.Lines())
		}
	}
	return
}

// FIX: this has to go into the std library
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
	//
	if c.Cmd("run rule", "print plural name").Begin() {
		// FIX no can do -- was trying to turn target into thing
		// what you need is something like:
		// c.Cmd("new", "print name", c.Cmd("get", "@", "target"))
		// where new treats everything
		c.Param("decide").Cmd("print text", c.Cmd("pluralize", c.Cmd("buffer", c.Cmd("determine", c.Cmd("get", "@", "target")))))
		c.End()
	}
	if c.Cmd("run rule", "print plural name").Begin() {
		c.Param("if").Cmd("is not", c.Cmd("is empty", c.Cmd("get", c.Cmd("get", "@", "target"), "printed plural name")))
		c.Param("decide").Cmd("print text", c.Cmd("get", c.Cmd("get", "@", "target"), "printed name"))
		c.End()
	}
}
