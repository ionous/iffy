package define_test

import (
	"github.com/ionous/iffy/dl/define"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/spec/ops"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestDefines(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var reg define.Registry
		var facts define.Facts
		e := reg.Define(&facts)
		testify.NoError(t, e)
	})
	//
	// t.Run("grammar", func(t *testing.T) {
	{
		var reg define.Registry
		reg.Register(defineGrammar)
		var facts define.Facts
		e := reg.Define(&facts)
		testify.NoError(t, e)
		//
		if testify.Len(t, facts.Grammar.Match, 1) {
			x, ok := facts.Grammar.Match[0].(*parser.AllOf)
			testify.True(t, ok)
			testify.Len(t, x.Match, 2) // l/look;action
		}
		// })
	}
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
func defineRelatives(*ops.Builder) {
}
func defineEventHandler(*ops.Builder) {
}
func definePattern(*ops.Builder) {
}
func defineRule(*ops.Builder) {
}

// future: defineClass, defineInstance, defineRelation, defineEvent,
