package evtbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
)

//
type Listeners struct {
	event.EventMap
	actions Actions
}

type ListenOn interface {
	On(string, event.Options, rt.Execute) error
}

func NewListeners(actions Actions) *Listeners {
	return &Listeners{make(event.EventMap), actions}
}

func (l *Listeners) Object(obj rt.Object) ListenOn {
	return phaseOn{Listeners: l, obj: obj, cls: obj.GetClass()}
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
	id := id.MakeId(name)
	// verify the action is well known.
	if _, ok := p.actions[id]; !ok {
		err = errutil.New("unknown action", name)
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
	}
	return
}
