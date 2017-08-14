package parser

import (
	"github.com/ionous/errutil"
	"strings"
)

// Multi matches one or more objects in ctx.
// (plus or minus some ambiguity)
type Multi struct {
	Filters Filters
}

func (try *Multi) Scan(ctx Context, cs Cursor) (ret Result, err error) {
	if word := cs.CurrentWord(); len(word) == 0 {
		err = MissingObject{Depth(cs.Pos)}
	} else {
		all := strings.EqualFold("all", word)
		if all {
			cs = cs.Skip(1)
		}

		r := &RankAll{Context: ctx, Filters: try.Filters}
		if !RankNouns(ctx.GetScope(), cs, r) {
			err = errutil.New("unexpected error")
		} else if !all && len(r.Plurals) == 0 {
			// we didnt have "all" and we didnt have plurals,
			// acts just as if one object was wanted.
			ret, err = resolveObject(cs, r.WordCount, r.Ranking.Nouns)
		} else {
			var nouns []Noun
			if r.Ranking.Empty() {
				nouns = r.Implied
			} else {
				nouns = r.Ranking.Nouns
			}

			if cnt := len(nouns); cnt > 0 {
				// filter nouns by any plurals
			MatchPlurals:
				for _, pl := range r.Plurals {
					for i := 0; i < cnt; {
						n := nouns[i]
						if n.HasPlural(pl) {
							i++
						} else if last := cnt - 1; last == 0 {
							// this unmatched object was the only thing in the list?
							// stop trying to filter objects from the list.
							break MatchPlurals
						} else {
							// slice out the one that didnt match
							nouns[i] = nouns[last]
							nouns, cnt = nouns[:last], last
						}
					}
				}
			}

			if len(nouns) == 0 {
				err = NoSuchObjects{Depth(cs.Pos + r.WordCount)}
			} else {
				wordCount := r.WordCount
				if all {
					wordCount++
				}
				ret = ResolvedMulti{nouns, wordCount}
			}
		}
	}
	return
}
