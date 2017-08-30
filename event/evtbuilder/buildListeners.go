package evtbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
)

//
type Listeners struct {
	eventTypes unique.Types
	event.EventMap
}

type ListenOn interface {
	On(string, event.Options, rt.Execute) error
}

func NewListeners(eventTypes unique.Types) *Listeners {
	return &Listeners{eventTypes, make(event.EventMap)}
}

func (l *Listeners) Object(obj rt.Object) ListenOn {
	return phaseOn{Listeners: l, obj: obj, cls: obj.Type()}
}

func (l *Listeners) Class(cls rt.Class) ListenOn {
	return phaseOn{Listeners: l, cls: cls}
}

type errOn struct {
	err error
}

func (errOn errOn) On(string, event.Options, rt.Execute) error {
	return errOn.err
}

type phaseOn struct {
	*Listeners
	obj rt.Object
	cls rt.Class
}

func (p phaseOn) On(name string, flags event.Options, exec rt.Execute) (err error) {
	id := ident.IdOf(name)
	if _, ok := p.eventTypes[id]; !ok {
		err = errutil.New("unknown event", name)
	} else {
		var target ident.Id
		var phaseMap event.PhaseMap
		listeners := p.EventMap[id]
		if p.obj != nil {
			target = p.obj.Id()
			if listeners.Objects != nil {
				phaseMap = listeners.Objects
			} else {
				phaseMap = make(event.PhaseMap)
				listeners.Objects = phaseMap
			}
		} else {
			target = class.Id(p.cls)
			if listeners.Classes != nil {
				phaseMap = listeners.Classes
			} else {
				phaseMap = make(event.PhaseMap)
				listeners.Classes = phaseMap
			}
		}
		//
		listenerType := event.BubbleListeners
		if flags.IsCapture() {
			listenerType = event.CaptureListeners
		}
		pl := phaseMap[target]
		pl[listenerType] = append(pl[listenerType], event.Handler{
			flags,
			exec,
		})
		phaseMap[target] = pl
		p.EventMap[id] = listeners
	}
	return
}
