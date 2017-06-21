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

//go:generate stringer -type=SecretState
type SecretState int

const (
	Secret SecretState = iota
	NotSecret
)

type DeadEnd struct {
	Secret SecretState
}

// TriState provides an enum with three choices for testing.
//go:generate stringer -type=TriState
type TriState int

const (
	No TriState = iota
	Yes
	Maybe
)

// TooLongState simulates an enum with an infinite number of values.
type TooLongState int

func (i TooLongState) String() string {
	return "repeats"
}

// EmptyState provides an enum with choices, but without stringer.
type EmptyState int

const (
	NotEmpty EmptyState = iota
)
