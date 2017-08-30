package std_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/locate"
	. "github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
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

	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*core.Classes)(nil),
		(*Classes)(nil))

	unique.RegisterBlocks(unique.PanicTypes(patterns),
		(*Patterns)(nil))

	objects := ref.NewObjects()
	unique.RegisterValues(unique.PanicValues(objects),
		Thingaverse.objects(sliceOf.String("box", "cake", "apple", "pen"))...)

	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*core.Commands)(nil),
		(*rule.Commands)(nil),
		(*Commands)(nil),
	)
	unique.RegisterTypes(unique.PanicTypes(cmds),
		(*Location)(nil),
	)

	// fix? if runtime was a set of slots, we could add a slot specifically for locale.
	assert := testify.New(t)
	rules, e := rule.Master(cmds, patterns, PrintNameRules, PrintObjectRules)
	assert.NoError(e)

	type OpsCb func(c *ops.Builder)
	type Match func(run rt.Runtime, lines []string) bool
	test := func(t *testing.T, build, exec OpsCb, match Match) (err error) {
		relations := ref.NewRelations()
		pc := locate.Locale{index.NewTable(index.OneToMany)}
		relations.AddTable("locale", pc.Table)

		var facts struct{ Locations []Locate }
		if c, ok := cmds.NewBuilder(&facts); !ok {
			err = errutil.New("no builder")
		} else {
			if c.Cmds().Begin() {
				build(c)
				c.End()
			}
			if e := c.Build(); e != nil {
				err = e
			} else {
				objs := objects.Build()

				for _, l := range facts.Locations {
					l := l.(*Location)
					// in this case we're probably a command too
					if p, ok := objs.GetObject(l.Parent); !ok {
						err = errutil.New("unknown", l.Parent)
						break
					} else if c, ok := objs.GetObject(l.Child); !ok {
						err = errutil.New("unknown", l.Child)
						break
					} else if e := pc.SetLocation(p, c, l.Locale); e != nil {
						err = e
						break
					}
				}
			}
		}
		if err == nil {
			var root struct{ rt.ExecuteList }
			if c, ok := cmds.NewBuilder(&root); !ok {
				err = errutil.New("no builder")
			} else {
				if c.Cmds().Begin() {
					exec(c)
					c.End()
				}
				if e := c.Build(); e != nil {
					err = e
				} else {
					var lines printer.Lines
					run := rtm.New(classes).Objects(objects).Rules(rules).Relations(relations).Writer(&lines).Rtm()
					if e := root.Execute(run); e != nil {
						err = e
					} else if res := lines.Lines(); match(run, res) {
						t.Logf("%s success: '%s'", t.Name(), strings.Join(res, ";"))
					}
				}
			}
		}
		return
	}
	//
	t.Run("contains", func(t *testing.T) {
		assert := testify.New(t)
		e := test(t, func(c *ops.Builder) {
			c.Cmd("Location", "box", locate.Contains, "cake")
		}, func(c *ops.Builder) {
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

	emptyBox := func(c *ops.Builder) {
	}
	boxContents := func(c *ops.Builder) {
		c.Cmd("Location", "box", locate.Contains, "cake")
		c.Cmd("Location", "box", locate.Contains, "apple")
		c.Cmd("Location", "box", locate.Contains, "pen")
	}
	printContent := func(c *ops.Builder) {
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
	printSummary := func(c *ops.Builder) {
		if c.Cmd("print span").Begin() {
			if c.Cmds().Begin() {
				c.Cmd("determine", c.Cmd("print summary", "box"))
				c.End()
			}
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
		return func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				c.Cmds(c.Cmd("determine", c.Cmd("print object", name)))
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
