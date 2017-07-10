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
// perhaps those have "default actions", but there's a difference b/t game actions --
// "default actions" here are simply functions.
// they have no event flow associated with them.
func (d *Dispatch) Go(run rt.Runtime, act *Action, target, data rt.Object) (err error) {
	if els, ok := d.events[act.Id]; ok {
		if path, e := els.CollectAncestors(run, target); e != nil {
			err = e
		} else {
			at := els.CollectTargets(target, nil)
			if len(path) > 0 || len(at) > 0 {
				evt := &EventObject{
					Name:          act.Name,
					Data:          data.GetId(),
					Bubbles:       true,
					Cancelable:    true,
					Target:        target.GetId(),
					CurrentTarget: target.GetId(),
				}
				// FIX: class finder and hints
				ac, e := NewFrame(run, evt)
				if e != nil {
					err = e
				} else if e := ac.DispatchFrame(at, path); e != nil {
					err = e
				} else if !evt.PreventDefault {
					evt.Phase = AfterPhase
					// FIX: set hint to target class
					if e := rt.ExecuteList(run, act.DefaultActions); e != nil {
						err = e
					} else {
						err = ac.Flush()
					}
				}
				// cleanup regardless of error
				ac.Destroy()
			}
		}
	}
	return
}
