package event

//go:generate stringer -type=Options
type Options int

const (
	Capture Options = (1 << iota)
	TargetOnly
	RunAfter
	Default Options = 0
)

func (opt Options) IsCapture() bool {
	return opt&Capture == Capture
}
func (opt Options) IsTargetOnly() bool {
	return opt&TargetOnly == TargetOnly
}
func (opt Options) IsRunAfter() bool {
	return opt&RunAfter == RunAfter
}

// EventPhase describes the event lifecycle.
type EventPhase int

//go:generate stringer -type=EventPhase
const (
	PhaseNone EventPhase = iota
	CapturingPhase
	AtTarget
	BubblingPhase
	AfterPhase
)

// ListenerType separates listeners into two categories.
// This is an implementation artifact.
type ListenerType int

//go:generate stringer -type=ListenerType
const (
	CaptureListeners ListenerType = iota
	BubbleListeners
	ListenerTypes
)
