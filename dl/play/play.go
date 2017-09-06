package play

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/event/evtbuilder"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"io"
)

var globalPlay Play

type Play struct {
	callbacks []func(*ops.Builder)
	classes   []interface{}
	cmds      []interface{}
	patterns  []interface{}
	events    []interface{}
	objects   []interface{}
}

// Register definitions globally. Used mainly via go init()
func (r *Play) AddScript(cb func(c *ops.Builder)) {
	r.callbacks = append(r.callbacks, cb)
}

// Register definitions globally. Used mainly via go init()
// func Register(cb func(c *ops.Builder)) {
// 	globalPlay.Register(cb)
// }

func (r *Play) AddClasses(block interface{}) {
	r.classes = append(r.classes, block)
}

func (r *Play) AddCommands(block interface{}) {
	r.cmds = append(r.cmds, block)
}

func (r *Play) AddPatterns(block interface{}) {
	r.patterns = append(r.patterns, block)
}

func (r *Play) AddEvents(block interface{}) {
	r.events = append(r.events, block)
}

func (r *Play) AddObjects(objs ...interface{}) {
	r.objects = append(r.objects, objs...)
}

// Define implements Statement by using all Register(ed) definitions.
func (r *Play) Build(cmds *ops.Ops) (ret Facts, err error) {
	var f Facts
	var root struct{ Definitions }
	if c, ok := cmds.NewBuilder(&root); ok {
		if c.Cmds().Begin() {
			for _, v := range r.callbacks {
				v(c)
			}
			c.End()
		}
		if e := c.Build(); e != nil {
			err = e
		} else if e := root.Define(&f); e != nil {
			err = e
		} else {
			ret = f
		}
	}
	return
}

func (r *Play) Play(w io.Writer) (ret *rtm.Rtm, err error) {
	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOpsX(classes, core.Xform{})    // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types
	events := unique.NewStack(patterns)           // all events become default action patterns
	objects := ref.NewObjects()
	relations := ref.NewRelations()

	//
	if e := unique.RegisterBlocks(classes, (*Classes)(nil)); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(classes, r.classes...); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(cmds, (*Commands)(nil)); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(cmds, r.cmds...); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(patterns, (*Patterns)(nil)); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(patterns, r.patterns...); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(events, r.events...); e != nil {
		err = e
	} else if e := unique.RegisterValues(objects, r.objects...); e != nil {
		err = e
	} else if rules, e := rule.Master(cmds, patterns, std.StdRules); e != nil {
		err = e
	} else if facts, e := r.Build(cmds); e != nil {
		err = e
	} else if e := facts.Mandates.Mandate(patterns.Types, rules); e != nil {
		err = e
	} else {
		// FIX: create a parser with facts.Grammar
		// noting, that we dont really have a parser yet -- just some teets.
		listen := evtbuilder.NewListeners(events.Types)

		pc := locate.Locale{index.NewTable(index.OneToMany)}
		if e := relations.AddTable("ParentChild", pc.Table); e != nil {
			err = e
		} else {
			run := rtm.New(classes).
				Objects(objects).
				Rules(rules).
				Ancestors(ParentChildAncestry{}).
				Relations(relations).
				Grammar(&facts.Grammar).
				Events(listen.EventMap).
				Writer(w).
				Rtm()

			if e := addLocations(run.ObjectMap, pc, facts.Locations); e != nil {
				err = e
			} else {
				if e := addObjectListeners(run.ObjectMap, listen, facts.ObjectListeners); e != nil {
					err = e
				} else if e := addClassListeners(run.Types, listen, facts.ClassListeners); e != nil {
					err = e
				} else {
					ret = run
				}
			}
		}
	}
	return
}

func addLocations(objs ref.ObjectMap, pc locate.Locale, ls []Location) (err error) {
	for _, loc := range ls {
		// in this case we're probably a command too
		if p, ok := objs.GetObject(loc.Parent); !ok {
			err = errutil.New("unknown", loc.Parent)
			break
		} else if c, ok := objs.GetObject(loc.Child); !ok {
			err = errutil.New("unknown", loc.Child)
			break
		} else if e := pc.SetLocation(p, c, loc.Locale); e != nil {
			err = e
			break
		}
	}
	return
}

func addObjectListeners(objs ref.ObjectMap, listen *evtbuilder.Listeners, ls []ListenTo) (err error) {
	for _, l := range ls {
		opt := l.GetOptions()
		if obj, ok := objs.GetObject(l.Target); !ok {
			err = errutil.New("couldnt find object", l.Target)
			break
		} else if e := listen.Object(obj).On(l.Event, opt, l.Go); e != nil {
			err = e
			break
		}
	}
	return
}

func addClassListeners(classes unique.Types, listen *evtbuilder.Listeners, ls []ListenFor) (err error) {
	for _, l := range ls {
		opt := l.GetOptions()
		// FIX: change to class registry
		if cls, ok := classes.FindType(l.Target); !ok {
			err = errutil.New("couldnt find class", l.Target)
			break
		} else if e := listen.Class(cls).On(l.Event, opt, l.Go); e != nil {
			err = e
			break
		}
	}
	return
}
