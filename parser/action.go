package parser

// Action terminates a matcher sequence, resolving to the named action.
type Action struct {
	Name string
}

// Scan matches only if the cursor has finished with all words.
func (a *Action) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if w := cs.CurrentWord(); len(w) == 0 {
		ret = ResolvedAction{a.Name}
	} else {
		err = Overflow{Depth(cs.Pos)}
	}
	return
}
