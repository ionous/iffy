package event

import (
	"github.com/ionous/iffy/rt"
)

type QueuedAction struct {
	src  rt.Object
	exec rt.Execute
}

type QueuedActions []QueuedAction

func (qs QueuedActions) Add(obj rt.Object, exec rt.Execute) QueuedActions {
	return append(qs, QueuedAction{obj, exec})
}

func (qs QueuedActions) Flush(run rt.Runtime, evt *EventObject) (err error) {
	// FIX? do we want the event object in scope during after handlers?
	// for consistancy's sake probably so.
	for _, qa := range qs {
		// FIX: set hint via hint pointer into scope
		evt.CurrentTarget = qa.src
		if e := qa.exec.Execute(run); e != nil {
			err = e
			break
		}
	}
	return
}
