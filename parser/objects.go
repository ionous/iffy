package parser

import (
	"github.com/ionous/errutil"
)

// Noun matches one object in ctx.
// (plus or minus some ambiguity)
type Noun struct {
	Filters Filters
}

func (try *Noun) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if w := cs.CurrentWord(); len(w) == 0 {
		err = MissingObject{Depth(cs.Pos)}
	} else {
		r := &RankOne{Filters: try.Filters}
		if !RankNouns(bounds, cs, r) {
			err = errutil.New("unexpected error")
		} else {
			ret, err = resolveObject(cs, r.Rank, r.Nouns)
		}
	}
	return
}

func resolveObject(cs Cursor, wordCount int, nouns []NounInstance) (ret Result, err error) {
	if wordCount == 0 {
		err = UnknownObject{Depth(cs.Pos)}
	} else {
		last := cs.Pos + wordCount
		if cnt := len(nouns); cnt == 1 {
			words := cs.Words[cs.Pos:last]
			ret = ResolvedObject{nouns[0], words}
		} else {
			err = AmbiguousObject{nouns, Depth(last)}
		}
	}
	return
}
