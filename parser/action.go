package parser

// Action terminates a matcher sequence by setting the context to the desired action.
type Action struct {
	Name string
}

// Scan matches only if the cursor has finished with all words.
func (a *Action) Scan(scope Scope, cs Cursor) (ret Result, err error) {
	if _, ok := cs.CurrentWord(); !ok {
		ret = ResolvedAction{a.Name}
	} else {
		err = Overflow{Depth(cs.Pos)}
	}
	return
}
