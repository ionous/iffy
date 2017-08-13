package parser

import (
	"strings"
)

// Object matches one or more objects in scope.
// (plus or minus some ambiguity)
type Multi struct {
	Filters Filters
}

func (try *Multi) Scan(scope Scope, cs Cursor) (ret Result, err error) {
	if word, ok := cs.CurrentWord(); !ok {
		err = MissingObject{Depth(cs.Pos)}
	} else if !strings.EqualFold("all", word) {
		// not all, then some one object
		ret, err = scanForObject(cs, scope, try.Filters)
	} else {
		cs := cs.Skip(1) // skip "all"
		if words, ids := matchObjects(cs, scope, try.Filters, true); words == 0 {
			ret = ResolvedMulti{Ids: ids}
		} else {
			words := cs.Words[cs.Pos : cs.Pos+words]
			ret = ResolvedMulti{ids, words}
		}
	}
	return
}
