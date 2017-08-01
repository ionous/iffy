package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/initial"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/pat/patspec"
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

func TestContents(t *testing.T) {
	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*core.Classes)(nil),
		(*Classes)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		Thingaverse.objects(sliceOf.String("box", "cake", "apple", "pen"))...)

	cmds := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*core.Commands)(nil),
		(*patspec.Commands)(nil),
		(*Commands)(nil),
		(*initial.Commands)(nil),
	)
	unique.RegisterBlocks(unique.PanicTypes(cmds.ShadowTypes),
		(*Patterns)(nil),
	)

	// fix? if runtime was a set of slots, we could add a slot specifically for locale.
	assert := testify.New(t)
	patterns, e := patbuilder.NewPatternMaster(cmds, classes,
		(*Patterns)(nil)).Build(printNamePatterns, printObjectPatterns)
	assert.NoError(e)

	type OpsCb func(c *ops.Builder)
	type Match func(run rt.Runtime, lines []string) bool
	test := func(t *testing.T, build, exec OpsCb, match Match) (err error) {
		relations := ref.NewRelations()
		pc := locate.Locale{index.NewTable(index.OneToMany)}
		relations.AddTable("locale", pc.Table)

		var src struct{ initial.Statements }
		if c, ok := cmds.NewBuilder(&src); !ok {
			err = errutil.New("no builder")
		} else {
			if c.Cmds().Begin() {
				build(c)
				c.End()
			}
			if e := c.Build(); e != nil {
				err = e
			} else {
				var facts initial.Facts
				if e := src.Assess(&facts); e != nil {
					err = e
				} else {
					objs := objects.Build()
					// for _, v := range facts.Values {
					// 	if obj, ok := objs.GetObject(v.Obj); !ok {
					// 		t.Fatal("couldnt find", v.Obj)
					// 		break
					// 	} else if e := obj.SetValue(v.Prop, v.Val); e != nil {
					// 		t.Fatal(e)
					// 		break
					// 	}
					// }
					// for _, r := range facts.Relations {
					// 	if e := relations.NewRelation(r.Name, index.NewTable(r.Type)); e != nil {
					// 		t.Fatal(e)
					// 		break
					// 	}
					// }

					for _, l := range facts.Locations {
						// in this case we're probably a command too
						if p, ok := objs.GetObject(l.Parent); !ok {
							err = errutil.New("unknown", l.Parent)
							break
						} else if c, ok := objs.GetObject(l.Child); !ok {
							err = errutil.New("unknown", l.Child)
							break
						} else if e := pc.SetLocation(p, c, l.Relative); e != nil {
							err = e
							break
						}
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
					run := rtm.New(classes).Objects(objects).Patterns(patterns).Relations(relations).Writer(&lines).Rtm()
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
	//
	emptyBox := func(c *ops.Builder) {
	}
	boxContents := func(c *ops.Builder) {
		c.Cmd("Location", "box", locate.Contains, "cake")
		c.Cmd("Location", "box", locate.Contains, "apple")
		c.Cmd("Location", "box", locate.Contains, "pen")
	}
	printContent := func(c *ops.Builder) {
		c.Cmd("determine", c.Cmd("print content", "box"))
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
			return assert.EqualValues(sliceOf.String("empire apple, cake, and pen"), lines)
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
	t.Run("with summary", func(t *testing.T) {
		assert := testify.New(t)
		Thingaverse["box"].(*Container).Closed = true
		e := test(t, emptyBox, printObject("box"), func(run rt.Runtime, lines []string) bool {
			return assert.EqualValues(sliceOf.String("box ( closed )"), lines)
		})
		assert.NoError(e)
	})
}
