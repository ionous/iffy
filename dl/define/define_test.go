package define_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/define"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/spec/ops"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestDefines(t *testing.T) {
	errutil.Panic = true
	t.Run("empty", func(t *testing.T) {
		var reg define.Registry
		//
		var facts define.Facts
		e := reg.Define(&facts)
		testify.NoError(t, e)
	})
	{
		var reg define.Registry
		reg.Register(defineGrammar)
		//
		var facts define.Facts
		e := reg.Define(&facts)
		testify.NoError(t, e)
		if testify.Len(t, facts.Grammar.Match, 1) {
			x, ok := facts.Grammar.Match[0].(*parser.AllOf)
			testify.True(t, ok)
			testify.Len(t, x.Match, 2) // l/look;action
		}
	}
	{
		var reg define.Registry
		reg.Register(defineLocation)
		//
		var facts define.Facts
		e := reg.Define(&facts)
		testify.NoError(t, e)
		testify.Len(t, facts.Locations, 1)
	}
	// {
	// 	var reg define.Registry
	// 	reg.Register(defineRules)
	// 	//
	// 	var facts define.Facts
	// 	e := reg.Define(&facts)
	// 	testify.NoError(t, e)
	// 	testify.Len(t, facts.Locations, 1)
	// }
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
				if c.Cmd("any of").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("action", "look")
						c.End()
					}
					c.End()
				}
				c.End()
			}
			c.End()
		}
		c.End()
	}
}
func defineLocation(c *ops.Builder) {
	c.Cmd("location", "parent", c.Cmd("supports"), "child")
}

func defineRules(c *ops.Builder) {
	for _, k := range []string{"bool", "number", "text", "object", "num list", "text list", "obj list", "run"} {
		if c.Cmd(k + " rule").Begin() {
			c.Param("name").Val(k)
			c.End()
		}
	}
}

func defineEventHandler(*ops.Builder) {

}

// future: defineEvent, defineInstance, defineRelative, defineClass, definePattern, defineRelation.
