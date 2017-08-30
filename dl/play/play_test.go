package play_test

import (
	// "github.com/ionous/iffy/dl/play"
	// "github.com/ionous/iffy/parser"
	// "github.com/ionous/iffy/spec/ops"
	// testify "github.com/stretchr/testify/assert"
	"testing"
)

// basic goal:
// everything is event based,
// the main game loop is a parser->event pump.
// so, a complete play test would be something like:
// parse "jump" and have the "player" jump
// what is a player? great question.
// Player.Actor as an object.
// define the gramar, and the handler in cmds
// dont worry about defineing the objects/classes as cmds.
func TestPlay(t *testing.T) {
	// 	var reg play.Registry
	// 	reg.Register(defineEventHandler)
	// 	//

	// 	type Jumping struct {
	// 	}
	// 	actions := evtbuilder.NewActions(classes, cmds)
	// 	actions.Add("jump", "kind", (*Jumping)(nil))
	// 	// FIX: add a registry filter so we can ensure the actions/events exist
	// 	listen := evtbuilder.NewListeners()

	// 	//
	// 	var facts play.Facts
	// 	e := reg.Define(&facts)
	// 	testify.NoError(t, e)
	// 	// could we build this directly?
	// 	// possibly pass an interface into listener
	// 	for _, lt := range facts.Listeners {
	// 		// fix: why do we have to determine this up front?
	// 		if GetObject(lt.Target) {
	// 			var opt event.Options
	// 			for _, p := range lt.Options {
	// 				opt |= p.Options()
	// 			}
	// 			// this takes a function, but we already  have executed it.
	// 			listen.Object(lt.Target).On(v.Event, opt, lt.Go)
	// 		}
	// 	}

	// 	//1:  we want to pass exec list.
	// 	if exec, e := fn.Build(p.actions.cmds); e != nil {
	// 		err = e
	// 	}

	// 	// no: we want to use patterns.
	// 	// actions.On("jump", func(c *ops.Builder) {
	// 	// 	if c.Cmd("print span").Begin() {
	// 	// 		if c.Cmds().Begin() {
	// 	// 			// FIX: to print names need to include articles
	// 	// 			// probably want a simple named object in core.
	// 	// 			c.Cmd("print text", "jumped!")
	// 	// 		}
	// 	// 	}
	// 	// })
}
