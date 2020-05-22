package chart

type State interface {
	StateName() string
	NewRune(rune) State
}
