package chart

type State interface {
	NewRune(rune) State
}

// Statement defines a function which satisfies the State interface.
type Statement func(rune) State

// NewRune implements State by calling the Statement's underlying function.
func (s Statement) NewRune(r rune) State { return s(r) }

// SelfStatement defines a function which satisfies the State interface by passing its function pointer as the first argument to the state callback. This allows closures to return themselves.
// For example the following definition:
//   recursive := func(self SelfStatement, r rune, pos Pos) (State) { return self }
// means "recursive" can return itself forever.
type SelfStatement func(SelfStatement, rune) State

// NewRune implements State by calling the underlying function.
func (s SelfStatement) NewRune(r rune) State { return s(s, r) }
