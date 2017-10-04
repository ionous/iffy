package std_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/dl/rules"
	. "github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/rel"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type Locate interface {
	Locate() Locate
}
type Location struct {
	Parent string
	Locale locate.Containment
	Child  string
}

func (l Location) Locate() Locate {
	return l
}

func TestContents(t *testing.T) {
	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOps(classes)                   // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types

	unique.PanicBlocks(classes,
		(*core.Classes)(nil),
		(*Classes)(nil))

	unique.PanicBlocks(patterns,
		(*Patterns)(nil))

	var objects obj.Registry
	objects.RegisterValues(Thingaverse.objects(
		sliceOf.String("box", "cake", "apple", "pen")))

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*rules.Commands)(nil),
		(*Commands)(nil),
	)
	unique.PanicTypes(cmds,
		(*Location)(nil),
	)

	// fix? if runtime was a set of slots, we could add a slot specifically for locale.
	assert := testify.New(t)
	rules, e := rules.Master(cmds, core.Xform{}, patterns, PrintNameRules, PrintObjectRules)
	assert.NoError(e)

	type OpsCb func(c spec.Block)
	type Match func(run rt.Runtime, lines []string) bool
	test := func(t *testing.T, build, exec OpsCb, match Match) (err error) {
		relations := rel.NewRelations()
		pc := locate.Locale{index.NewTable(index.OneToMany)}
		relations.AddTable("locale", pc.Table)

		var loc struct{ Locations []Locate }
		c := cmds.NewBuilder(&loc, core.Xform{})
		if e := c.Build(build); e != nil {
			err = e
		} else {
			var root struct{ rt.ExecuteList }
			c := cmds.NewBuilder(&root, core.Xform{})
			if e := c.Build(exec); e != nil {
				err = e
			} else {
				var lines printer.Lines
				if run, e := rtm.New(classes).Objects(objects).Rules(rules).Relations(relations).Writer(&lines).Rtm(); e != nil {
					for _, l := range loc.Locations {
						l := l.(*Location)
						// in this case we're probably a command too
						if p, ok := run.GetObject(l.Parent); !ok {
							err = errutil.New("unknown", l.Parent)
							break
						} else if c, ok := run.GetObject(l.Child); !ok {
							err = errutil.New("unknown", l.Child)
							break
						} else if e := pc.SetLocation(p, l.Locale, c); e != nil {
							err = e
							break
						}
					}

					if err == nil {
						if e := root.Execute(run); e != nil {
							err = e
						} else if res := lines.Lines(); match(run, res) {
							t.Logf("%s success: '%s'", t.Name(), strings.Join(res, ";"))
						}
					}
				}
			}
		}
		return
	}
	//
	t.Run("contains", func(t *testing.T) {
		assert := testify.New(t)
		e := test(t, func(c spec.Block) {
			c.Cmd("Location", "box", locate.Contains, "cake")
		}, func(c spec.Block) {
			//
		}, func(run rt.Runtime, lines []string) (okay bool) {
			if pc, ok := run.GetRelation("locale"); assert.True(ok) {
				if in, ok := pc.GetTable().GetData("$box", "$cake"); assert.True(ok) {
					okay = assert.EqualValues(locate.Contains, in)
				}
			}
			return
		})
		assert.NoError(e)
	})

	emptyBox := func(c spec.Block) {
	}
	boxContents := func(c spec.Block) {
		c.Cmd("Location", "box", locate.Contains, "cake")
		c.Cmd("Location", "box", locate.Contains, "apple")
		c.Cmd("Location", "box", locate.Contains, "pen")
	}
	printContent := func(c spec.Block) {
		c.Cmd("determine", c.Cmd("print content", "box", c.Param("tersely").Val(true)))
	}
	t.Run("empty", func(t *testing.T) {
		assert := testify.New(t)
		e := test(t, emptyBox, printContent, func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("empty"), lines)
		})
		assert.NoError(e)
	})
	t.Run("has contents", func(t *testing.T) {
		assert := testify.New(t)
		e := test(t, boxContents, printContent, func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("empire apple", "cake", "pen"), lines)
		})
		assert.NoError(e)
	})
	// summary tests:
	printSummary := func(c spec.Block) {
		if c.Cmd("print span").Begin() {
			c.Cmd("determine", c.Cmd("print summary", "box"))
			c.End()
		}
	}
	t.Run("closed", func(t *testing.T) {
		assert := testify.New(t)
		Thingaverse["box"].(*Container).Closed = true
		e := test(t, boxContents, printSummary, func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("closed"), lines)
		})
		assert.NoError(e)
	})
	t.Run("open but empty", func(t *testing.T) {
		assert := testify.New(t)
		Thingaverse["box"].(*Container).Closed = false
		e := test(t, emptyBox, printSummary, func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("open but empty"), lines)
		})
		assert.NoError(e)
	})
	t.Run("open contents", func(t *testing.T) {
		assert := testify.New(t)
		Thingaverse["box"].(*Container).Closed = false
		e := test(t, boxContents, printSummary, func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("in which is an empire apple, a cake, and a pen"), lines)
		})
		assert.NoError(e)
	})
	// print object: simple name, name with summary ( for a container )
	printObject := func(name string) OpsCb {
		return func(c spec.Block) {
			if c.Cmd("print span").Begin() {
				c.Cmd("determine", c.Cmd("print object", name))
				c.End()
			}
		}
	}
	t.Run("without summary", func(t *testing.T) {
		assert := testify.New(t)
		e := test(t, emptyBox, printObject("pen"), func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("pen"), lines)
		})
		assert.NoError(e)
	})
	t.Run("with closed summary", func(t *testing.T) {
		assert := testify.New(t)
		Thingaverse["box"].(*Container).Closed = true
		e := test(t, emptyBox, printObject("box"), func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("box ( closed )"), lines)
		})
		assert.NoError(e)
	})
	t.Run("with open summary", func(t *testing.T) {
		assert := testify.New(t)
		Thingaverse["box"].(*Container).Closed = false
		e := test(t, boxContents, printObject("box"), func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("box ( in which is an empire apple, a cake, and a pen )"), lines)
		})
		assert.NoError(e)
	})
}
