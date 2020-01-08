package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
)

func MissingKinds(db *sql.DB, cb func(string)) error {
	return missingNames(db,
		`select distinct name from named
		where not exists (
			select 1 from ancestry a
			where named.name == a.kind
			and named.category = 'kind'
		)`, cb)
}

func MissingFields(db *sql.DB, cb func(string)) error {
	return missingNames(db,
		`select distinct name from named
		where not exists (
			select 1 from property p
			where named.name == p.field
			and named.category = 'field'
		)`, cb)
}

func missingNames(db *sql.DB, q string, cb func(string)) (err error) {
	if kinds, e := db.Query(q); e != nil {
		err = e
	} else {
		for kinds.Next() {
			var k string
			if e := kinds.Scan(&k); e != nil {
				err = e
				break
			} else {
				cb(k)
			}
		}
		if e := kinds.Err(); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
