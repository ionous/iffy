package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/initial"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/pat/patspec"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

// start looking at printing, probably a match after saying something.
// two blocks possibly: one for the execute match, one for the statements
// no "events" yet.

type Locale struct {
	*index.Table
}

func (l *Locale) SetLocation(p, c rt.Object, now locate.Containment) (err error) {
	types := map[locate.Containment]struct {
		Parent, Child string
	}{
		locate.Supports: {Parent: "container", Child: "thing"},
		locate.Contains: {Parent: "container", Child: "thing"},
		locate.Wears:    {Parent: "actor", Child: "thing"},
		locate.Carries:  {Parent: "actor", Child: "thing"},
		locate.Holds:    {Parent: "room", Child: "thing"},
	}
	if check, ok := types[now]; !ok {
		err = errutil.New("relation not supported", now)
	} else if !p.GetClass().IsCompatible(check.Parent) {
		err = errutil.New("expected parent", check.Parent)
	} else if !c.GetClass().IsCompatible(check.Child) {
		err = errutil.New("expected child", check.Child)
	} else {
		err = l.Table.AddPair(p.GetId(), c.GetId(), func(old interface{}) (ret interface{}, err error) {
			if c, ok := old.(locate.Containment); ok && c != now {
				err = errutil.New("was", c, "now", now)
			} else {
				ret = now
			}
			return
		})
	}
	return
}

func TestContents(t *testing.T) {
	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		Thingaverse.objects(sliceOf.String("box", "cake"))...)

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

	//t.Run("Names", func(t *testing.T) {
	assert := testify.New(t)

	// patterns, e := patbuilder.NewPatternMaster(cmds, classes,
	// 	(*Patterns)(nil)).Build(printPatterns)
	// assert.NoError(e)

	pc := Locale{index.NewTable(index.OneToMany)}

	run := func(c *ops.Builder) {
		c.Cmd("Location", "box", locate.Contains, "cake")
	}

	var root struct{ initial.Statements }
	if c, ok := cmds.NewBuilder(&root); ok {
		if c.Cmds().Begin() {
			run(c)
			c.End()
		}
		if e := c.Build(); assert.NoError(e) {
			var facts initial.Facts
			if e := root.Assess(&facts); assert.NoError(e) {
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
						t.Fatal(l.Parent)
						break
					} else if c, ok := objs.GetObject(l.Child); !ok {
						t.Fatal(l.Child)
						break
					} else if e := pc.SetLocation(p, c, l.Relative); e != nil {
						t.Fatal(e)
					}
				}
				// verify relation:
				if in, ok := pc.GetData("$box", "$cake"); assert.True(ok) {
					assert.EqualValues(locate.Contains, in)
				}
			}

		}
	}
	// })
}
