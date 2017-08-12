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
type Ranking struct {
	Rank int
	Ids  []string
}

func (try *Object) Scan(scope Scope, scan Cursor) (ret Result, err error) {
	if !scan.HasNext() {
		err = MissingObject{Depth(scan.Pos)}
	} else {
		var best Ranking
		scope.SearchScope(func(n Noun) bool {
			if matchName(n, scan) && try.Filters.MatchesNoun(n) {
				cs := scan.Skip(1)
				// determine the ranking of this object by eating words as long as they match.
				var rank int
				for ; matchName(n, cs); rank++ {
					cs = cs.Skip(1)
				}
				switch id := n.GetId(); {
				case rank > best.Rank:
					best = Ranking{rank, sliceOf.String(id)}
				case rank == best.Rank:
					best.Ids = append(best.Ids, id)
				}
			}
			return false // always keep going
		})

		ids := best.Ids
		if cnt := len(ids); cnt == 0 {
			err = UnknownObject{Depth(scan.Pos)}
		} else {
			last := scan.Pos + best.Rank + 1
			if cnt == 1 {
				words := scan.Words[scan.Pos:last]
				ret = ResolvedObject{ids[0], words}
			} else {
				err = AmbiguousObject{ids, Depth(last)}
			}
		}
	}
	return
}

func matchName(n Noun, cs Cursor) bool {
	name, ok := cs.GetNext()
	return ok && n.HasName(name)
}
