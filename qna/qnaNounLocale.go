package qna

import (
	"database/sql"
)

type nounLocale struct {
	q     *sql.Stmt         // relativesOf
	cache map[string]string // child to parent
}

func (a *nounLocale) reset() {
	// its okay to clear this because setting locale records to the db.
	a.cache = nil
}

func (a *nounLocale) setLocaleOf(id, parent string) {
	if a.cache == nil {
		a.cache = make(map[string]string)
	}
	a.cache[id] = parent
}

func (a *nounLocale) localeOf(id string) (ret string) {
	if el, ok := a.cache[id]; ok {
		ret = el
	} else {
		switch e := a.q.QueryRow(id, "locale").Scan(&el); e {
		default:
			panic(e)
		case nil, sql.ErrNoRows:
			a.setLocaleOf(id, el)
			ret = el
		}
	}
	return
}
