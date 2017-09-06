package event_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/event/evtbuilder"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
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
		Jumper rt.Object `if:"cls:kind"`
	}

	type Kiss struct {
		Kisser   rt.Object `if:"cls:kind"`
		KissWhom rt.Object `if:"cls:kind"`
	}

	type Unlock struct {
		Unlocker rt.Object `if:"cls:kind"`
		Lock     rt.Object `if:"cls:kind"`
		With     rt.Object `if:"cls:kind"`
	}

	type Events struct {
		*Jump
		*Kiss
		*Unlock
	}

	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOpsX(classes, core.Xform{})    // all shadow types become classes
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
						c.Cmd("say", "jumped!")
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
			c.Cmd("say", "bogart's jumping!")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.Default, jump)
			assert.NoError(e)
		}
	}
	{
		if jump, e := cmds.Execute(func(c *ops.Builder) {
			c.Cmd("say", "bogart's going to jump!")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.Capture, jump)
			assert.NoError(e)
		}
	}
	{
		if jump, e := cmds.Execute(func(c *ops.Builder) {
			c.Cmd("say", "bogart's tired of jumping.")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.RunAfter, jump)
			assert.NoError(e)
		}
	}

	jump := run.Emplace(&Jump{Jumper: bogart})
	if obj, e := event.TargetOf(run, jump); assert.NoError(e) {
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
			c.Cmd("say", "don't do it bogart!")
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

}
