package parser

import (
	"github.com/ionous/errutil"
)

// Focus scanner provides a way to change scope for subsequent scanners.
// For instance, searching only though held objects.
type Focus struct {
	// future: Who string -- with "" means¥ing player
	What  Scanner
	Where string
}

//
func (a *Focus) Scan(ctx Context, _ Scope, cs Cursor) (ret Result, err error) {
	if scope, e := ctx.GetPlayerScope(a.Where); e != nil {
		err = e
	} else {
		ret, err = a.What.Scan(ctx, scope, cs)
	}
	return
}

// Target changes the scope of its first scanner in response to the results of its last scanner. Generally, this means that the last scanner should be Noun{}.
type Target struct {
	Match []Scanner
}

//
func (a *Target) Scan(ctx Context, scope Scope, start Cursor) (ret Result, err error) {
	first, rest := a.Match[0], AllOf{a.Match[1:]}
	errorDepth := -1
	// scan ahead for matches and determine how many words might match this target.
	for cs := start; len(cs.CurrentWord()) > 0; cs = cs.Skip(1) {
		if rl, e := rest.scan(ctx, scope, cs); e != nil {
			// like any of, we track the "deepest" error.
			if d := DepthOf(e); d > errorDepth {
				err, errorDepth = e, d
			}
			continue // keep looking for success
		} else if last, ok := rl.Last(); !ok {
			err = errutil.New("target not found")
		} else if obj, ok := last.(ResolvedObject); !ok {
			err = errutil.Fmt("expected an object, got %T", last)
			break
		} else if scope, e := ctx.GetOtherScope(obj.NounVisitor.GetId()); e != nil {
			err = e
			break
		} else {
			// snip down our input to just the preface which matches this target
			// ( to avoid errors of "too many words" )
			words := start.Words[:cs.Pos]
			sub := Cursor{start.Pos, words}
			if r, e := first.Scan(ctx, scope, sub); e != nil {
				err = e
				break
			} else {
				rl.AddResult(r)
				ret, err = rl, nil
				break
			}
		}
	}
	return
}
