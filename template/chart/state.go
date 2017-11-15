package chart

type State interface {
	NewRune(rune) State
}

// Statement functions behave as a State.
type Statement func(rune) State

// NewRune implements State by calling the Statement's underlying function.
func (s Statement) NewRune(r rune) State { return s(r) }

// SelfStatement are State functions which pass the function pointer as the first argument to the state callback. This allows closures to return themselves.
// For example, the following state returns itself forever:
//   var recursive  SelfStatement = func(self SelfStatement, r rune) State { return self }
type SelfStatement func(SelfStatement, rune) State

// NewRune implements State by calling the underlying function.
func (s SelfStatement) NewRune(r rune) State { return s(s, r) }

// Action are states which always return nil.
// see also: StateExit.
type Action func()

// NewRune implements State by calling the action and returning nil.
func (s Action) NewRune(r rune) State {
	s()
	return nil
}

// StateExit marks a terminating state for readability's sake.
func StateExit(action Action) State {
	return action
}
