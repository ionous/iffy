package event

import (
	"github.com/ionous/iffy/rt"
)

type Dispatch struct {
	events EventMap
}

func NewDispatch(events EventMap) *Dispatch {
	return &Dispatch{events}
}

// note: in the future, other things may raise events.
// for instance, object state change.
func (d *Dispatch) Go(run rt.Runtime, act *Action, target, data rt.Object) (err error) {
	if els, ok := d.events[act.Id]; ok {
		if path, e := els.CollectAncestors(run, target); e != nil {
			err = e
		} else {
			at := els.CollectTargets(target, nil)
			if len(path) > 0 || len(at) > 0 {
				evt := &EventObject{
					Name:          act.Name,
					Data:          data,
					Bubbles:       true,
					Cancelable:    true,
					Target:        target,
					CurrentTarget: target,
				}
				// FIX: class finder and hints
				ac, e := NewFrame(run, evt)
				if e != nil {
					err = e
				} else if e := ac.DispatchFrame(at, path); e != nil {
					err = e
				} else if !evt.PreventDefault {
					evt.Phase = AfterPhase
					evt.CurrentTarget = target
					// FIX: set hint to target class
					if e := act.DefaultActions.Execute(run); e != nil {
						err = e
					} else {
						err = ac.queue.Flush(run, evt)
					}
				}
				// cleanup regardless of error
				ac.Destroy()
			}
		}
	}
	return
}
