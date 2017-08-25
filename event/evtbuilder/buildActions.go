package evtbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
)

type Actions event.ActionMap

// Add registers a new action builder class.
func (a Actions) Add(name string, targetClass, dataClass rt.Class) (err error) {
	id := id.MakeId(name)
	if act, exists := a[id]; exists {
		err = errutil.New("duplicate action registered", name, act.Name)
	} else {
		a[id] = &event.Action{
			id,
			name,
			targetClass,
			dataClass,
			nil,
		}
	}
	return
}

// On adds a default action.
func (a Actions) On(name string, exec rt.Execute) (err error) {
	id := id.MakeId(name)
	if act, exists := a[id]; !exists {
		err = errutil.New("unknown action", name)
	} else {
		act.DefaultActions = append(act.DefaultActions, exec)
	}
	return
}
