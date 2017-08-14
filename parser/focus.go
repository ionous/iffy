package parser

// Focus scanner provides a way to change scope for subsequent scanners.
// For instance, searching only though held objects.
type Focus struct {
	// future: Who string -- with "" means¥ing player
	What  Scanner
	Where string
}

//
func (a *Focus) Scan(ctx Context, _ Scope, cs Cursor) (ret Result, err error) {
	scope := ctx.GetPlayerScope(a.Where)
	return a.What.Scan(ctx, scope, cs)
}
