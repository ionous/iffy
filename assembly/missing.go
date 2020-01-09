package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
)

// MissingKinds returns named kinds which don't have a defined ancestry.
func MissingKinds(db *sql.DB, cb func(string)) error {
	var k string
	return queryAll(db,
		`select distinct name from named
		where not exists (
			select 1 from ancestry a
			where named.name == a.kind
			and named.category = 'kind'
		)`, func() error { cb(k); return nil }, &k)
}

// MissingFields returns named fields which don't have a defined property.
func MissingFields(db *sql.DB, cb func(string)) error {
	var k string
	return queryAll(db,
		`select distinct name from named
		where not exists (
			select 1 from property p
			where named.name == p.field
			and named.category = 'field'
		)`, func() error { cb(k); return nil }, &k)
}

func queryAll(db *sql.DB, q string, cb func() error, dest ...interface{}) (err error) {
	if it, e := db.Query(q); e != nil {
		err = e
	} else {
		for it.Next() {
			if e := it.Scan(dest...); e != nil {
				err = e
				break
			} else if e := cb(); e != nil {
				err = e
				break
			}
		}
		if e := it.Err(); e != nil {
			err = errutil.Append(err, e)
		}
		it.Close()
	}
	return
}
