package event

import (
	"github.com/ionous/iffy/rt"
)

// Map event id to listener list
type EventMap map[string]EventListeners

// we have two kinds of listeners, class listeners, and object listeners.
type EventListeners struct {
	Classes, Objects PhaseMap
}

type PhaseMap map[string]PhaseList

// PhaseList contains capture and bubble handlers.
type PhaseList [AfterPhase][]*Handler

type Handler struct {
	Options
	Object rt.Object
	Class  rt.Class
	Exec   rt.Execute
}

type Target struct {
	obj      rt.Object
	cls      rt.Class
	handlers PhaseList
}

// CollectAncestors to create targets from the parents of the pased object. The  target order is: instance's parent, parent classes, container instance, repeat.
func (els EventListeners) CollectAncestors(run rt.Runtime, obj rt.Object) (ret []Target, err error) {
	if at, e := run.GetAncestors(obj); e != nil {
		err = e
	} else {
		var tgt []Target
		for at.HasNext() {
			if obj, e := at.GetNext(); e != nil {
				err = e
				break
			} else {
				tgt = els.CollectTargets(obj, tgt)
			}
		}
		if err == nil {
			ret = tgt
		}
	}
	return
}

func (els EventListeners) CollectTargets(obj rt.Object, tgt []Target) []Target {
	// check instance listeners
	if ls, ok := els.Objects[obj.GetId()]; ok {
		tgt = append(tgt, Target{obj: obj, handlers: ls})
	}
	// check class listeners
	for cls := obj.GetClass(); ; {
		if ls, ok := els.Objects[cls.GetId()]; ok {
			tgt = append(tgt, Target{obj: obj, cls: cls, handlers: ls})
		}
		// move to parent class
		if next, ok := cls.GetParent(); !ok {
			break
		} else {
			cls = next
		}
	}
	return tgt
}
