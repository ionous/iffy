package parser

import (
	"github.com/ionous/sliceOf"
)

// Object matches one object in scope.
// (plus or minus some ambiguity)
type Object struct {
	Filters Filters
}

func (try *Object) Scan(scan Cursor, ctx *Context) (ret int, okay bool) {
	var best Ranking
	ctx.SearchScope(func(n Noun) bool {
		if matchName(n, scan, ctx) && try.Filters.MatchesNoun(n) {
			cs := scan + 1
			// keep eating words as long as they match this object.
			// FUTURE? reverse ambiguity, allow the next match to grow backwards if it can, increasing the ambiguity of the phrase
			// ex. look inside: did you mean look to the inside, or look inside something...
			var rank int
			for ; matchName(n, cs, ctx); rank++ {
				cs++
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
	return
}

func matchName(n Noun, cs Cursor, ctx *Context) bool {
	name, ok := cs.NextWord(ctx)
	return ok && n.HasName(name)
}
