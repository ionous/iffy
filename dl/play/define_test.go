package play

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/spec/ops"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

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

// Define implements Statement by using all AddScript(ed) definitions.
func (r *Play) Define(f *Facts) (err error) {
	classes := make(unique.Types)
	unique.RegisterBlocks(
		unique.PanicTypes(classes),
		(*Classes)(nil),
	)

	cmds := ops.NewOpsX(classes, core.Xform{})
	unique.RegisterBlocks(
		unique.PanicTypes(cmds),
		(*Commands)(nil),
	)

	unique.RegisterBlocks(
		unique.PanicTypes(cmds.ShadowTypes),
		(*Patterns)(nil),
	)

	var root struct{ Definitions }
	if c, ok := cmds.NewBuilder(&root); ok {
		if c.Cmds().Begin() {
			for _, v := range r.callbacks {
				v(c)
			}
			c.End()
		}
		if e := c.Build(); e != nil {
			err = e
		} else {
			err = root.Define(f)
		}
	}
	return
}

func TestEmpty(t *testing.T) {
	var reg Play
	var facts Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
}

func TestGrammar(t *testing.T) {
	var reg Play
	reg.AddScript(defineGrammar)
	//
	var facts Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	if testify.Len(t, facts.Grammar.Match, 1) {
		x, ok := facts.Grammar.Match[0].(*parser.AllOf)
		testify.True(t, ok)
		testify.Len(t, x.Match, 2) // l/look;action
	}
}

func TestLocation(t *testing.T) {
	var reg Play
	reg.AddScript(defineLocation)
	//
	var facts Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	testify.Len(t, facts.Locations, 1)
}

func TestRules(t *testing.T) {
	var reg Play
	mandates := []string{"bool", "number", "text", "object", "num list", "text list", "obj list", "run"}
	reg.AddScript(func(c *ops.Builder) {
		defineRules(c, mandates)
	})
	//
	var facts Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	testify.Len(t, facts.Mandates, len(mandates))
}

func TestEvents(t *testing.T) {
	var reg Play
	reg.AddScript(defineEventHandler)
	//
	var facts Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	testify.Len(t, facts.ObjectListeners, 1)
}

func defineGrammar(c *ops.Builder) {
	if c.Cmd("grammar").Begin() {
		if c.Cmd("all of").Begin() {
			if c.Cmds().Begin() {
				if c.Cmd("any of").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("word", "l")
						c.Cmd("word", "look")
						c.End()
					}
					c.End()
				}
				c.Cmd("action", "look")
				c.End()
			}
			c.End()
		}
		c.End()
	}
}

func defineLocation(c *ops.Builder) {
	c.Cmd("location", "parent", locate.Supports, "child")
}

func defineRules(c *ops.Builder, mandates []string) {
	for _, k := range mandates {
		if c.Cmd("mandate").Begin() {
			if c.Cmd(k + " rule").Begin() {
				c.Param("name").Val(k)
				c.End()
			}
			c.End()
		}
	}
}

func defineEventHandler(c *ops.Builder) {
	if c.Cmd("listen to", "bogart", "jump").Begin() {
		if c.Param("go").Cmds().Begin() {
			c.Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "target")))
			c.Cmd("print text", "jumping!")
			c.End()
		}
		if c.Param("options").Cmds().Begin() {
			c.Cmd("capture")
			c.Cmd("target only")
			c.End()
		}
		c.End()
	}
	if c.Cmd("mandate").Begin() {
		if c.Cmd("run rule", "jump").Begin() {
			if c.Param("decide").Cmds().Begin() {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "target")))
						c.Cmd("print text", "jumped!")
						c.End()
					}
					c.End()
				}
				c.End()
			}
			c.Param("continue").Cmd("continue after")
			c.End()
		}
		c.End()
	}
}
