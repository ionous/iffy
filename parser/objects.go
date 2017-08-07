package parser

import (
	"github.com/ionous/sliceOf"
	// "strings"
)

// Object matches one object in scope.
// (plus or minus some ambiguity)
type Object struct {
	Filters Filters
}

func (try *Object) Scan(scan Cursor, ctx *Context) (ret int, okay bool) {
	var best Ranking
	if !scan.HasNext() {
		// the player hasnt typed anything else, but we are the last noun
		// ask them which object they mean.
		ctx.NeedsNoun = true
		okay = true
	} else {
		ctx.SearchScope(func(n Noun) bool {
			if matchName(n, scan) && try.Filters.MatchesNoun(n) {
				cs := scan.Step(1)
				// keep eating words as long as they match this object.
				// FUTURE? reverse ambiguity, allow the next match to grow backwards if it can, increasing the ambiguity of the phrase
				// ex. look inside: did you mean look to the inside, or look inside something...
				var rank int
				for ; matchName(n, cs); rank++ {
					cs = cs.Step(1)
				}
				switch id := n.GetId(); {
				case rank > best.Rank:
					best = Ranking{rank, sliceOf.String(id)}
				case rank == best.Rank:
					best.Nouns = append(best.Nouns, id)
				}
			}
			return false // keep going
		})
		if len(best.Nouns) > 0 {
			// FIX: can we make this a return?
			ctx.Results.Matches = append(ctx.Results.Matches, best)
			ret, okay = best.Rank+1, true // number of words
		}
	}
	return
}

func matchName(n Noun, cs Cursor) bool {
	name, ok := cs.GetNext()
	return ok && n.HasName(name)
}
