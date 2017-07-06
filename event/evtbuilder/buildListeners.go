package evtbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
)

//
type Listeners struct {
	event.EventMap
	actions *Actions
	objects *ref.Objects
	classes *ref.Classes
}

type ListenOn interface {
	On(name string, flags event.Options, fn BuildOps) error
}

func NewListeners(actions *Actions, objects *ref.Objects, classes *ref.Classes) *Listeners {
	return &Listeners{make(event.EventMap), actions, objects, classes}
}

func (l *Listeners) Object(name string) (ret ListenOn) {
	if obj, ok := l.objects.GetObject(name); !ok {
		ret = errOn{errutil.New("unknown object", name)}
	} else {
		ret = phaseOn{Listeners: l, obj: obj, cls: obj.GetClass()}
	}
	return
}

func (l *Listeners) Class(name string) (ret ListenOn) {
	if cls, ok := l.classes.GetClass(name); !ok {
		ret = errOn{errutil.New("unknown object", name)}
	} else {
		ret = phaseOn{Listeners: l, cls: cls}
	}
	return
}

type errOn struct {
	err error
}

func (errOn errOn) On(string, event.Options, BuildOps) error {
	return errOn.err
}

type phaseOn struct {
	*Listeners
	obj rt.Object
	cls rt.Class
}

func (p phaseOn) On(name string, flags event.Options, fn BuildOps) (err error) {
	id := id.MakeId(name)
	// verify the action is well known.
	if _, ok := p.actions.ActionMap[id]; !ok {
		err = errutil.New("unknown action", name)
	} else if exec, e := fn.Build(p.actions.ops); e != nil {
		err = e
	} else {
		var target string
		var phaseMap event.PhaseMap
		// use the action id as our event id
		el := p.EventMap[id]
		if p.obj != nil {
			target = p.obj.GetId()
			if el.Objects != nil {
				phaseMap = el.Objects
			} else {
				phaseMap = make(event.PhaseMap)
				p.EventMap[id] = el
				el.Objects = phaseMap
			}
		} else {
			target = p.cls.GetId()
			if el.Classes != nil {
				phaseMap = el.Classes
			} else {
				phaseMap = make(event.PhaseMap)
				p.EventMap[id] = el
				el.Classes = phaseMap
			}
		}
		//
		phase := event.BubblePhase
		if flags.IsCapture() {
			phase = event.CapturePhase
		}
		pl := phaseMap[target]
		pl[phase] = append(pl[phase], &event.Handler{
			flags,
			p.obj,
			p.cls,
			exec,
		})
		phaseMap[target] = pl
	}
	return
}
