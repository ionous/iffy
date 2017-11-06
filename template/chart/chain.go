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

// NewRune tries the first state, and any of its returned states; then switches to the last state.
func (p *ChainParser) NewRune(r rune) (ret State) {
	if next := p.next.NewRune(r); next != nil {
		ret, p.next = p, next
	} else {
		ret = p.last.NewRune(r)
	}
	return
}

// MakeParallel region; run all of the passed states until they all return nil.
func MakeParallel(rs ...State) State {
	return SelfStatement(func(self SelfStatement, r rune) (ret State) {
		var cnt int
		for _, s := range rs {
			if next := s.NewRune(r); next != nil {
				rs[cnt] = next
				cnt++
			}
		}
		if cnt > 0 {
			rs = rs[:cnt]
			ret = self
		}
		return
	})
}

// for the very next rune returns nil ( unhandled )
var Terminal = Statement(func(rune) State { return nil })
