package chart

// parseChain sends passed rune to the chain of first and last states.
func parseChain(r rune, first, last State) State {
	return makeChain(first, last).NewRune(r)
}

// makeChain connects two states so that if a rune is not handled by the first state, the rune is delegated to the second state.
// this is similar to a parent-child statechart relationship.
func makeChain(first, last State) State {
	return Statement(func(r rune) (ret State) {
		if next := first.NewRune(r); next != nil {
			ret = makeChain(next, last)
		} else {
			ret = last.NewRune(r)
		}
		return
	})
}
