package qna

import (
	"database/sql"
)

type activeNouns struct {
	q     *sql.Stmt
	cache map[string]bool
}

func (a *activeNouns) reset() {
	a.cache = nil
}

func (a *activeNouns) isActive(id string) (ret bool) {
	if el, ok := a.cache[id]; ok {
		ret = el
	} else {
		switch e := a.q.QueryRow(id).Scan(&el); e {
		default:
			panic(e)
		case nil, sql.ErrNoRows:
			if a.cache == nil {
				a.cache = make(map[string]bool)
			}
			a.cache[id] = el
			ret = el
		}
	}
	return
}
