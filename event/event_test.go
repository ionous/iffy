package event_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/event/evtbuilder"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
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

	cmds := ops.NewOps((*core.Commands)(nil))
	classes := ref.NewClasses()
	objects := ref.NewObjects(classes)
	relations := ref.NewRelations(classes, objects)

	unique.RegisterTypes(
		unique.PanicTypes(classes),
		(*Kind)(nil))

	unique.RegisterValues(
		unique.PanicValues(objects),
		&Kind{"Bogart"},
		&Kind{"Bob"},
		&Kind{"Coffin"},
		&Kind{"Skeleton Key"})

	// this isnt parser
	actions := evtbuilder.NewActions(classes, cmds)
	actions.Add("jump", "kind", (*Jumping)(nil))
	actions.Add("kiss", "kind", (*Kissing)(nil))
	actions.Add("unlock", "kind", (*Unlocking)(nil))

	// default action:
	actions.On("jump", func(c *ops.Builder) {
		if c.Cmd("print line").Begin() {
			if c.Cmds().Begin() {
				// FIX: to print names need to include articles
				// probably want a simple named object in core.
				c.Cmd("print text", "jumped!")
			}
		}
	})

	// we do this the manual way first, and later with spec
	listen := evtbuilder.NewListeners(actions, objects, classes)
	// object listener:
	listen.Object("bogart").On("jump", event.Default, func(c *ops.Builder) {
		c.Cmd("print text", "bogart jumped!")
	})
	listen.Class("kind").On("kiss", event.Default, func(c *ops.Builder) {
		c.Cmd("print text", "kissed!")
	})

	var lines rtm.LineWriter
	run := rtm.NewRtm(classes, objects, relations)
	run.PushWriter(&lines)

	dispatch := event.NewDispatch(listen.EventMap)

	// helper for testing:
	Go := func(object, action string, dataFn evtbuilder.BuildOps) {
		if obj, ok := objects.GetObject(object); !ok {
			t.Fatal("object not found", object)
		} else if act, ok := actions.ActionMap[id.MakeId(action)]; !ok {
			t.Fatal("unknown action", action)
		} else {
			var data rt.Object
			if dataFn != nil {
				if dataEval, e := dataFn.Eval(cmds); e != nil {
					t.Fatal(e)
				} else if got, e := dataEval.GetObject(run); e != nil {
					t.Fatal(e)
				} else {
					data = got
				}
			}
			e := dispatch.Go(run, act, obj, data)
			assert.NoError(e)
		}
	}

	Go("bogart", "jump", nil)
	t.Log(lines.Lines())

	// Go("bogart", "kiss", func(c *ops.Builder) {
	// 	c.Value("bob")
	// })
	// Go("bob", "unlock", func(c *ops.Builder) {
	// 	c.Param("lock", "coffin")
	// 	c.Param("with", "skeleton key")
	// })

}
