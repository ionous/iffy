package parser

// Action terminates a matcher sequence by setting the context to the desired action.
type Action struct {
	Name string
}

// Scan matches only if the cursor has finished with all words.
func (a *Action) Scan(scope Scope, scan Cursor) (ret Result, err error) {
	if !scan.HasNext() {
		ret = ResolvedAction{a.Name}
	} else {
		err = Overflow{Depth(scan.Pos)}
	}
	return
}
