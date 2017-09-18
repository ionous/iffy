package play

import (
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/event/trigger"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/spec"
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
	play.AddObjects(
		&std.Story{Name: "story"},
		&std.Room{Kind: std.Kind{Name: "Circus"}},
		&std.Actor{Thing: std.Thing{Kind: std.Kind{Name: "Bogart",
			CommonProper: std.ProperNamed}}},
		&std.Pawn{"pawn", ident.IdOf("bogart")},
	)

	var lines printer.Lines
	if run, e := play.Play(&lines); e != nil {
		t.Fatal(e)
	} else {
		// 4. initial position
		// 3. banner
		// 2. evaluate status bars
		// 1. commence game
		// 5. "read" input
		// 6. trigger related command
		// 7. end turn - updates status bar. bc what if the score changed during? inital: 0/1
		// When play begins: say "Welcome to Old Marston Grange, a country house cut off by fog."
		// ie. verify output order is roughly like inform7
		// 8. start adding commands -- test them separately, and try a little story which uses the here.
		if e := statup(run); e != nil {
			t.Fatal(e)
		} else {
			if e := trigger.Parse(run, "jump"); e != nil {
				t.Fatal(e)
			} else {
				assert.Equal(sliceOf.String("Circus", "0/0", "Bogart is jumping!", "Bogart jumped!"), lines.Lines())
			}
		}
	}
}

// FIX: verify at compile time that all the actions mentioned by grammar exist.
func definePlay(c spec.Block) {
	if c.Cmd("grammar").Begin() {
		if c.Cmd("all of").Begin() {
			if c.Cmds().Begin() {
				c.Cmd("word", "jump")
				if c.Cmd("trigger").Begin() {
					c.Cmd("jump", c.Cmd("player"))
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
	if c.Cmd("location").Begin() {
		c.Val("circus").Val(locate.Has).Val("bogart")
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
