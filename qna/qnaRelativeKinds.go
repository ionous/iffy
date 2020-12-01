package qna

import (
	"database/sql"
	"strings"
)

type relativeKinds struct {
	q     *sql.Stmt
	cache map[string]relativeKind
}

type relativeKind struct {
	// kind and other kind are simple names
	// noun kind paths can then be tested against them with compatibleKind
	// ex. if a relation supports "animals",
	// a noun of "kinds,animals,cats" can be used with it.
	//
	kind, otherKind, cardinality string
}

// doesThis, implementThat
func compatibleKind(path, k string) bool {
	return strings.Contains(path, k+",")
}

func (a *relativeKinds) relativeKind(id string) (ret relativeKind) {
	if el, ok := a.cache[id]; ok {
		ret = el
	} else {
		switch e := a.q.QueryRow(id).Scan(&el.kind, &el.otherKind, &el.cardinality); e {
		default:
			panic(e)
		case nil, sql.ErrNoRows:
			if a.cache == nil {
				a.cache = make(map[string]relativeKind)
			}
			a.cache[id] = el
			ret = el
		}
	}
	return
}
