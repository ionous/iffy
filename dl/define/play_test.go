package define_test

import (
	"github.com/ionous/iffy/dl/define"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/spec/ops"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestPlay(t *testing.T) {
	var reg define.Registry
	reg.Register(defineEventHandler)
	//

	type Jumping struct {
	}
	actions := evtbuilder.NewActions(classes, cmds)
	actions.Add("jump", "kind", (*Jumping)(nil))
	//
	listen := evtbuilder.NewListeners(actions, run.Objects, classes)

	//
	var facts define.Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	// could we build this directly?
	// possibly pass an interface into listener
	for _, lt := range facts.Listeners {
		// fix: why do we have to determine this up front?
		if GetObject(lt.Target) {
			var opt event.Options
			for _, p := range lt.Options {
				opt |= p.Options()
			}
			// this takes a function, but we already  have executed it.
			listen.Object(lt.Target).On(v.Event, opt, lt.Go)
		}
	}

	//1:  we want to pass exec list.
	if exec, e := fn.Build(p.actions.cmds); e != nil {
		err = e
	}

	// no: we want to use patterns.
	// actions.On("jump", func(c *ops.Builder) {
	// 	if c.Cmd("print span").Begin() {
	// 		if c.Cmds().Begin() {
	// 			// FIX: to print names need to include articles
	// 			// probably want a simple named object in core.
	// 			c.Cmd("print text", "jumped!")
	// 		}
	// 	}
	// })
}
