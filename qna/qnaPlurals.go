package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
)

func NewPlurals(db *sql.DB) (ret *Plurals, err error) {
	var ps tables.Prep
	n := &Plurals{
		singularOf: ps.Prep(db, `select one 
		from mdl_plural
		where many=?
		limit 1`),
		pluralOf: ps.Prep(db, `select many 
		from mdl_plural
		where one=?
		limit 1`),
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = n
	}
	return
}

type Plurals struct {
	singularOf, pluralOf *sql.Stmt
}

func (n *Plurals) Singular(str string) (ret string, err error) {
	return n.get(n.singularOf, str, lang.Singularize)
}

func (n *Plurals) Plural(str string) (ret string, err error) {
	return n.get(n.pluralOf, str, lang.Pluralize)
}

func (n *Plurals) get(s *sql.Stmt, str string,
	mani func(string) string) (ret string, err error) {
	if s == nil {
		err = errutil.New("invalid statement")
	} else {
		var res string
		switch e := s.QueryRow(str).Scan(&res); e {
		case nil:
			// res was scanned in successfully.
		case sql.ErrNoRows:
			res = mani(str)
		default:
			res, err = str, e
		}
		// could cache the results if desired
		// actually measuring the code and lookup first...
		ret = res
	}
	return
}
