package event_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/event/evtbuilder"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
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

	type Jumping struct {
	}

	type Kissing struct {
		KissWhom *Kind
	}

	type Unlocking struct {
		Lock *Kind
		With *Kind
	}

	cmds := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*core.Commands)(nil))

	classes := ref.NewClasses()
	objects := ref.NewObjects(classes)

	unique.RegisterTypes(
		unique.PanicTypes(classes),
		(*Kind)(nil))

	unique.RegisterValues(
		unique.PanicValues(objects),
		&Kind{"Bogart"},
		&Kind{"Bob"},
		&Kind{"Coffin"},
		&Kind{"Skeleton Key"})

	dataClasses := ref.NewClassStack(classes)
	unique.RegisterTypes(dataClasses,
		(*Jumping)(nil),
		(*Kissing)(nil),
		(*Unlocking)(nil),
	)

	// this isnt parser
	actions := make(evtbuilder.Actions)
	kind, ok := classes.GetClass("kind")
	assert.True(ok)

	jumping, _ := dataClasses.GetClass("jumping")
	actions.Add("jump", kind, jumping)

	kissing, _ := dataClasses.GetClass("kissing")
	actions.Add("kiss", kind, kissing)

	unlocking, _ := dataClasses.GetClass("unlocking")
	actions.Add("unlock", kind, unlocking)

	// default action:
	if jumped, e := cmds.Execute(func(c *ops.Builder) {
		if c.Cmd("print span").Begin() {
			if c.Cmds().Begin() {
				// FIX: to print names need to include articles
				// probably want a simple named object in core.
				c.Cmd("print text", "jumped!")
				c.End()
			}
			c.End()
		}
	}); e != nil {
		t.Fatal(e)
	} else {
		actions.On("jump", jumped)
	}

	// we do this the manual way first, and later with spec

	var lines printer.Lines
	run := rtm.New(classes).Objects(objects).Writer(&lines).Rtm()

	listen := evtbuilder.NewListeners(actions)
	// object listener:
	if jump, e := cmds.Execute(func(c *ops.Builder) {
		c.Cmd("print text", "bogart jumped!")
	}); e != nil {
		t.Fatal(e)
	} else {
		bogart, _ := run.GetObject("bogart")
		listen.Object(bogart).On("jump", event.Default, jump)
	}

	if kiss, e := cmds.Execute(func(c *ops.Builder) {
		c.Cmd("print text", "kissed!")
	}); e != nil {
		t.Fatal(e)
	} else {
		listen.Class(kind).On("kiss", event.Default, kiss)
	}

	dispatch := event.NewDispatch(listen.EventMap)

	// helper for testing:
	Go := func(object, action string) {
		if obj, ok := run.GetObject(object); !ok {
			t.Fatal("object not found", object)
		} else if act, ok := actions[id.MakeId(action)]; !ok {
			t.Fatal("unknown action", action)
		} else {
			var data rt.Object
			/*if dataFn != nil {
				if dataEval, e := dataFn.Eval(cmds); e != nil {
					t.Fatal(e)
				} else if got, e := dataEval.GetObject(run); e != nil {
					t.Fatal(e)
				} else {
					data = got
				}
			}*/
			e := dispatch.Go(run, act, obj, data)
			assert.NoError(e)
		}
	}

	Go("bogart", "jump")
	t.Log(lines.Lines())

	// Go("bogart", "kiss", func(c *ops.Builder) {
	// 	c.Value("bob")
	// })
	// Go("bob", "unlock", func(c *ops.Builder) {
	// 	c.Param("lock", "coffin")
	// 	c.Param("with", "skeleton key")
	// })

}
