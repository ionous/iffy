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

//go:generate stringer -type=Phase
type Phase int

const (
	BubblePhase Phase = iota
	CapturePhase
	AfterPhase
)
