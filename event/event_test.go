package event_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/event/evtbuilder"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec"
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
		Jumper ident.Id `if:"cls:kind"`
	}

	type Kiss struct {
		Kisser   ident.Id `if:"cls:kind"`
		KissWhom ident.Id `if:"cls:kind"`
	}

	type Unlock struct {
		Unlocker ident.Id `if:"cls:kind"`
		Lock     ident.Id `if:"cls:kind"`
		With     ident.Id `if:"cls:kind"`
	}

	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOps(classes)                   // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types
	events := unique.NewStack(patterns)           // all events become default action patterns

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*rules.Commands)(nil))

	unique.PanicTypes(classes,
		(*Kind)(nil))

	unique.RegisterTypes(events,
		(*Jump)(nil))

	objects := obj.NewObjects()
	unique.PanicValues(objects,
		&Kind{"Bogart"},
		&Kind{"Bob"},
		&Kind{"Coffin"},
		&Kind{"Skeleton Key"})

	// default action:
	DefaultActions := func(c spec.Block) {
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

	rules, e := rules.Master(cmds, core.Xform{}, patterns, DefaultActions)
	assert.NoError(e)

	// we do this the manual way first, and later with spec

	var lines printer.Lines
	run := rtm.New(classes).Objects(objects).Rules(rules).Writer(&lines).Rtm()

	listen := evtbuilder.NewListeners(events.Types)
	// object listener:
	bogart, _ := run.GetObject("bogart")
	{
		if jump, e := Execute(cmds, func(c spec.Block) {
			c.Cmd("say", "bogart's jumping!")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.Default, jump)
			assert.NoError(e)
		}
	}
	{
		if jump, e := Execute(cmds, func(c spec.Block) {
			c.Cmd("say", "bogart's going to jump!")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.Capture, jump)
			assert.NoError(e)
		}
	}
	{
		if jump, e := Execute(cmds, func(c spec.Block) {
			c.Cmd("say", "bogart's tired of jumping.")
		}); e != nil {
			t.Fatal(e)
		} else {
			e := listen.Object(bogart).On("jump", event.RunAfter, jump)
			assert.NoError(e)
		}
	}

	jump := run.Emplace(&Jump{Jumper: bogart.Id()})
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
		if jump, e := Execute(cmds, func(c spec.Block) {
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
	run.SetWriter(&lines)
	//
	assert.NoError(event.Trigger(run, listen.EventMap, jump))
	assert.Equal(sliceOf.String("don't do it bogart!"), lines.Lines())
}

func Execute(cmds *ops.Ops, fn func(c spec.Block)) (ret rt.Execute, err error) {
	var root struct{ Eval rt.ExecuteList }
	c := cmds.NewBuilder(&root, core.Xform{})
	if c.Cmds().Begin() {
		fn(c)
		c.End()
	}
	if e := c.Build(); e != nil {
		err = e
	} else {
		ret = root.Eval
	}
	return
}
