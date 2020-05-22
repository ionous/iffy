package chart

// Action are states which always return nil.
// see also: StateExit.
func Action(name string, closure func()) State {
	return &actionState{name, closure}
}

type actionState struct {
	name    string
	closure func()
}

func (s *actionState) StateName() string {
	return s.name
}

// NewRune implements State by calling the action and returning nil.
func (s *actionState) NewRune(r rune) State {
	s.closure()
	return nil
}

// StateExit marks a terminating state for readability's sake.
func StateExit(name string, onExit func()) State {
	return &actionState{"exit " + name, onExit}
}

// for the very next rune returns nil ( unhandled )
var Terminal = Statement("terminal", func(rune) State {
	return nil
})
