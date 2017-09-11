package play

import (
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/event/trigger"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

// basic goal:
// the main game loop is a parser->event pump.
// so, a complete play test would be something like:
// parse "jump" and have the "player" jump
// what is a player? great question.
// Player.Actor as an object.
// define the gramar, and the handler in cmds
// dont worry about defineing the objects/classes as cmds.
func TestPlay(t *testing.T) {
	assert := testify.New(t)

	type Jump struct {
		Jumper ident.Id `if:"cls:kind"`
	}
	type Events struct {
		*Jump
	}
	var play Play
	play.AddEvents((*Events)(nil))
	play.AddScript(definePlay)
	bogart := &std.Actor{Thing: std.Thing{Kind: std.Kind{Name: "Bogart", CommonProper: std.ProperNamed}}}
	play.AddObjects(bogart, &std.Player{Name: "player", Pawn: ident.IdOf("bogart")})
	// errutil.Panic = true

	//couldnt find $printName

	var lines printer.Lines
	if run, e := play.Play(&lines); assert.NoError(e) {
		// FIX: verify at compile time that all the actions mentioned by grammar exist.
		if e := trigger.Parse(run, "jump"); assert.NoError(e) {
			assert.Equal(sliceOf.String("Bogart is jumping!", "Bogart jumped!"), lines.Lines())
		}
	}
}

func definePlay(c *ops.Builder) {
	if c.Cmd("grammar").Begin() {
		if c.Cmd("all of").Begin() {
			if c.Cmds().Begin() {
				c.Cmd("word", "jump")
				if c.Cmd("trigger").Begin() {
					c.Cmd("jump", c.Cmd("get", "player", "pawn"))
					c.End()
				}
				c.End()
			}
			c.End()
		}
		c.End()
	}
	if c.Cmd("mandate").Begin() {
		if c.Cmd("run rule", "jump").Begin() {
			if c.Param("decide").Cmds().Begin() {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "jumper")))
						c.Cmd("say", "jumped!")
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
	if c.Cmd("listen to", "bogart", "jump").Begin() {
		if c.Param("go").Cmds().Begin() {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("upper the", c.Cmd("get", c.Cmd("get", "@", "data"), "jumper"))
					c.Cmd("say", "is jumping!")
					c.End()
				}
				c.End()
			}
			c.End()
		}
		c.End()
	}
}
