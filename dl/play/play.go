package play

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/express"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/event/evtbuilder"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/rel"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
	"io"
)

type Play struct {
	callbacks []func(spec.Block)
	classes   []interface{}
	cmds      []interface{}
	patterns  []interface{}
	events    []interface{}
	objects   obj.Registry
}

func (p *Play) AddClasses(block interface{}) {
	p.classes = append(p.classes, block)
}
func (p *Play) AddCommands(block interface{}) {
	p.cmds = append(p.cmds, block)
}
func (p *Play) AddEvents(block interface{}) {
	p.events = append(p.events, block)
}
func (p *Play) AddObjects(objs ...interface{}) {
	p.objects.RegisterValues(objs)
}
func (p *Play) AddPatterns(block interface{}) {
	p.patterns = append(p.patterns, block)
}
func (p *Play) AddScript(cb func(c spec.Block)) {
	p.callbacks = append(p.callbacks, cb)
}

// Define implements Statement by using all Register(ed) definitions.
func (p *Play) build(cmds *ops.Ops, xform ops.Transform) (ret Facts, err error) {
	var f Facts
	var root struct{ Definitions }
	c := cmds.NewBuilder(&root, xform)
	if e := c.Build(p.callbacks...); e != nil {
		err = e
	} else if e := root.Define(&f); e != nil {
		err = e
	} else {
		ret = f
	}
	return
}

func (p *Play) Play(w io.Writer) (ret *rtm.Rtm, err error) {
	classes := make(unique.Types)                 // all types known to iffy
	cmds := ops.NewOps(classes)                   // all shadow types become classes
	patterns := unique.NewStack(cmds.ShadowTypes) // all patterns are shadow types
	events := unique.NewStack(patterns)           // all events become default action patterns
	relations := rel.NewRelations()
	//
	if e := unique.RegisterBlocks(classes, (*Classes)(nil)); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(classes, p.classes...); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(cmds, (*Commands)(nil)); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(cmds, p.cmds...); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(patterns, (*Patterns)(nil)); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(patterns, p.patterns...); e != nil {
		err = e
	} else if e := unique.RegisterBlocks(events, p.events...); e != nil {
		err = e
	} else {
		xform := express.NewTransform(cmds, &p.objects)
		if rules, e := rules.Master(cmds, xform, patterns, std.Rules); e != nil {
			err = e
		} else if facts, e := p.build(cmds, xform); e != nil {
			err = e
		} else if e := facts.Mandates.Mandate(rules); e != nil {
			err = e
		} else {
			// FIX: create a parser with facts.Grammar
			// noting, that we dont really have a parser yet -- just some teets.
			listen := evtbuilder.NewListeners(events.Types)

			pc := locate.Locale{index.NewTable(index.OneToMany)}
			if e := relations.AddTable("locale", pc.Table); e != nil {
				err = e
			} else if run, e := rtm.New(classes).
				Objects(p.objects).
				Rules(rules).
				Ancestors(ParentChildAncestry{}).
				Relations(relations).
				Grammar(&facts.Grammar).
				Events(listen.EventMap).
				Writer(w).
				Rtm(); e != nil {
				err = e
			} else if e := addLocations(run.Objects, pc, facts.Locations); e != nil {
				err = e
			} else if e := addObjectListeners(run.Objects, listen, facts.ObjectListeners); e != nil {
				err = e
			} else if e := addClassListeners(run.Types, listen, facts.ClassListeners); e != nil {
				err = e
			} else {
				ret = run
			}
		}
	}
	return
}

func addLocations(objs obj.ObjectMap, pc locate.Locale, ls []Location) (err error) {
	for _, loc := range ls {
		// in this case we're probably a command too
		if p, ok := objs.GetObject(loc.Parent); !ok {
			err = errutil.New("unknown", loc.Parent)
			break
		} else if c, ok := objs.GetObject(loc.Child); !ok {
			err = errutil.New("unknown", loc.Child)
			break
		} else if e := pc.SetLocation(p, loc.Locale, c); e != nil {
			err = e
			break
		}
	}
	return
}

func addObjectListeners(objs obj.ObjectMap, listen *evtbuilder.Listeners, ls []ListenTo) (err error) {
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
