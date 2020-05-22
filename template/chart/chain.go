package chart

// ParseChain of states, sending the passed rune to the chain.
func ParseChain(r rune, first, last State) State {
	return MakeChain(first, last).NewRune(r)
}

// MakeChain of states by connecting two states.
// If a rune is not handled by the first state or any of its returned states,
// the rune is handed to the second state.
// This is similar to a parent-child statechart relationship.
func MakeChain(first, last State) State {
	return &ChainParser{first, last}
}

// ChainParser: see MakeChain.
type ChainParser struct {
	next, last State
}

func (p *ChainParser) StateName() string {
	return "chain parser ('" + p.next.StateName() + "' '" + p.last.StateName() + "')"
}

// NewRune tries the first state, and any of its returned states; then switches to the last state.
func (p *ChainParser) NewRune(r rune) (ret State) {
	if next := p.next.NewRune(r); next != nil {
		ret, p.next = p, next
	} else {
		ret = p.last.NewRune(r)
	}
	return
}
