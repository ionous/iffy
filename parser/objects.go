package parser

import (
	"github.com/ionous/sliceOf"
)

// Object matches one object in scope.
// (plus or minus some ambiguity)
type Object struct {
	Filters Filters
}

func (try *Object) Scan(ctx *Context, start *Cursor) (ret int) {
	var best Ranking
	ctx.SearchScope(func(n Noun) bool {
		cs := *start
		if matchName(n, &cs) && try.Filters.MatchesNoun(n) {
			// keep eating words as long as they match this object.
			// FUTURE? reverse ambiguity, allow the next match to grow backwards if it can, increasing the ambiguity of the phrase
			// ex. look inside: did you mean look to the inside, or look inside something...
			var rank int
			for ; matchName(n, &cs); rank++ {
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
		ret = best.Rank + 1 // number of words
	}
	return
}

func matchName(n Noun, cs *Cursor) bool {
	name, ok := cs.NextWord()
	return ok && n.HasName(name)
}
