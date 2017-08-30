package event_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/event/evtbuilder"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/assert"
	"testing"
)

// this is a full fledged integration test
// would also need some smaller things.
func TestSomething(t *testing.T) {
	assert := assert.New(t)

	type Kind struct {
		Name string `if:"id"`
	}

	// FIX: add test to find target
	// ex. put it as the second object in a structure.
	type Jump struct {
		Jumper *Kind
	}

	type Kiss struct {
		Kisser   *Kind
		KissWhom *Kind
	}

	type Unlock struct {
		Unlocker *Kind
		Lock     *Kind
		With     *Kind
	}

	type Events struct {
		*Jump
		*Kiss
		*Unlock
	}

	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOps(classes)                   // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types
	events := unique.NewStack(patterns)           // all events become default action patterns

	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*core.Commands)(nil),
		(*rule.Commands)(nil))

	unique.RegisterTypes(
		unique.PanicTypes(classes),
		(*Kind)(nil))

	unique.RegisterBlocks(
		unique.PanicTypes(events),
		(*Events)(nil))

	objects := ref.NewObjects()
	unique.RegisterValues(
		unique.PanicValues(objects),
		&Kind{"Bogart"},
		&Kind{"Bob"},
		&Kind{"Coffin"},
		&Kind{"Skeleton Key"})

	// default action:
	DefaultActions := func(c *ops.Builder) {
		if c.Cmd("run rule", "jump").Begin() {
			if c.Param("decide").Cmds().Begin() {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						// FIX: to print names need to include articles
						// probably want a simple named object in core.
						c.Cmd("print text", "jumped!")
						c.End()
					}
					c.End()
				}
				c.End()
			}
			c.End()
		}
	}

	rules, e := rule.Master(cmds, patterns, DefaultActions)
	assert.NoError(e)

	// we do this the manual way first, and later with spec

	var lines printer.Lines
	run := rtm.New(classes).Objects(objects).Rules(rules).Writer(&lines).Rtm()

	listen := evtbuilder.NewListeners(events.Types)
	// object listener:
	bogart, _ := run.GetObject("bogart")
	{
		if jump, e := cmds.Execute(func(c *ops.Builder) {
			c.Cmd("print text", "bogart's jumping!")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.Default, jump)
			assert.NoError(e)
		}
	}
	{
		if jump, e := cmds.Execute(func(c *ops.Builder) {
			c.Cmd("print text", "bogart's going to jump!")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.Capture, jump)
			assert.NoError(e)
		}
	}
	{
		if jump, e := cmds.Execute(func(c *ops.Builder) {
			c.Cmd("print text", "bogart's tired of jumping.")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.RunAfter, jump)
			assert.NoError(e)
		}
	}

	// if kiss, e := cmds.Execute(func(c *ops.Builder) {
	// 	c.Cmd("print text", "kissed!")
	// }); e != nil {
	// 	t.Fatal(e)
	// } else {
	// 	listen.Class(kind).On("kiss", event.Default, kiss)
	// }

	// helper for testing:
	// Go := func(object, action string) {
	// 		if obj, ok := run.GetObject(object); !ok {
	// 			t.Fatal("object not found", object)
	// 		} else if act, ok := actions[ident.IdOf(action)]; !ok {
	// 			t.Fatal("unknown action", action)
	// 		} else {
	// 			var data rt.Object
	// 			if dataFn != nil {
	// 				if dataEval, e := dataFn.Eval(cmds); e != nil {
	// 					t.Fatal(e)
	// 				} else if got, e := dataEval.GetObject(run); e != nil {
	// 					t.Fatal(e)
	// 				} else {
	// 					data = got
	// 				}
	// 			}
	// 			e := dispatch.Go(run, act, obj, data)
	// 			assert.NoError(e)
	// 		}

	jump, e := run.Objects.Emplace(&Jump{
		Jumper: bogart.(*ref.RefObject).Value.Addr().Interface().(*Kind),
	})
	if obj, e := event.TargetOf(jump); assert.NoError(e) {
		assert.Equal(bogart, obj)
	}

	if els, ok := listen.EventMap[ident.IdOf("jump")]; assert.True(ok) {
		at := els.CollectTargets(bogart, nil)
		assert.Len(at, 1)
	}

	assert.NoError(event.Trigger(run, listen.EventMap, jump))
	assert.Equal(sliceOf.String("bogart's going to jump!", "bogart's jumping!", "jumped!", "bogart's tired of jumping."), lines.Lines())

	{
		if jump, e := cmds.Execute(func(c *ops.Builder) {
			c.Cmd("print text", "don't do it bogart!")
			c.Cmd("set bool", "@", "stop immediate propagation", true)
			c.Cmd("set bool", "@", "prevent default", true)
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.Capture, jump)
			assert.NoError(e)
		}
	}
	lines = printer.Lines{}
	run.Writer = &lines
	//
	assert.NoError(event.Trigger(run, listen.EventMap, jump))
	assert.Equal(sliceOf.String("don't do it bogart!"), lines.Lines())

	// Go("bogart", "kiss", func(c *ops.Builder) {
	// 	c.Value("bob")
	// })
	// Go("bob", "unlock", func(c *ops.Builder) {
	// 	c.Param("lock", "coffin")
	// 	c.Param("with", "skeleton key")
	// })

}
