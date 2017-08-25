package event

import (
	"github.com/ionous/iffy/rt"
)

type EventObject struct {
	Name string    // name of event
	Data rt.Object // data for the event

	Bubbles    bool // r/o does this bubble
	Cancelable bool // r/o can this event be canceled

	Target        rt.Object  // originator of event
	Phase         EventPhase // event flow phase
	CurrentTarget rt.Object  // current object processing the event

	PreventDefault           bool // stop the default action from running
	StopPropagation          bool // stop the event flow after the current target
	StopImmediatePropagation bool // stop processing all other event handlers immediately
}

func (evt *EventObject) Stopped() bool {
	return evt.StopPropagation || evt.StopImmediatePropagation
}
