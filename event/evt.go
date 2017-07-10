package event

type EventObject struct {
	Name              string // name of event
	Data/*Id*/ string // data for the event

	Bubbles    bool // r/o does this bubble
	Cancelable bool // r/o can this event be canceled

	Target/*Id*/ string        // originator of event
	Phase                      EventPhase // event flow phase
	CurrentTarget/*Id*/ string // current object processing the event

	PreventDefault           bool // stop the default action from running
	StopPropagation          bool // stop the event flow after the current target
	StopImmediatePropagation bool // stop processing all other event handlers immediately
}

func (evt *EventObject) Stopped() bool {
	return evt.StopPropagation || evt.StopImmediatePropagation
}
