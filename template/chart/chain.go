package chart

func parseChain(r rune, first, last State) State {
	return makeChain(first, last).NewRune(r)
}

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
