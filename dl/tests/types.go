package tests

//go:generate stringer -type=ListeningState
type ListeningState int

//go:generate stringer -type=StandingState
type StandingState int

const (
	Listening ListeningState = iota
	Laughing
)

const (
	Standing StandingState = iota
	Sitting
)

type Person struct {
	Name      string `if:"id"`
	Listening ListeningState
	Standing  StandingState
}
