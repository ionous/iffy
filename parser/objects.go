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

func (try *Object) Scan(scope Scope, cs Cursor) (ret Result, err error) {
	if _, ok := cs.CurrentWord(); !ok {
		err = MissingObject{Depth(cs.Pos)}
	} else {
		ret, err = scanForObject(cs, scope, try.Filters)
	}
	return
}

func scanForObject(cs Cursor, scope Scope, filters Filters) (ret Result, err error) {
	if words, ids := matchObjects(cs, scope, filters, false); words == 0 {
		err = UnknownObject{Depth(cs.Pos)}
	} else {
		last := cs.Pos + words
		if cnt := len(ids); cnt == 1 {
			words := cs.Words[cs.Pos:last]
			ret = ResolvedObject{ids[0], words}
		} else {
			err = AmbiguousObject{ids, Depth(last)}
		}
	}
	return
}

var matchAll = Filters{}

// MatchObjects returns a list of objects and the number of words which match them.
// FIX? for now, when ret is 0, ids is the whole scope.
func matchObjects(cs Cursor, scope Scope, filters Filters, all bool) (ret int, ids []string) {
	// visit every object in scope.
	scope.SearchScope(func(n Noun) bool {
		if filters.MatchesNoun(n) {
			if rank := RankNoun(cs, n); rank == 0 && all {
				ids = append(ids, n.GetId())
			} else if rank > 0 {
				all = false
				switch id := n.GetId(); {
				case rank > ret:
					ret, ids = rank, sliceOf.String(id)
				case rank == ret:
					ids = append(ids, id)
				}
			}
		}
		return false // never stop
	})
	return
}

// RankNoun returns how many words in src match the passed noun.
func RankNoun(src Cursor, n Noun) (rank int) {
	for ; matchName(n, src); rank++ {
		src = src.Skip(1)
	}
	return
}

func matchName(n Noun, cs Cursor) bool {
	name, ok := cs.CurrentWord()
	return ok && n.HasName(name)
}
