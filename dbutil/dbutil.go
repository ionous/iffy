package dbutil

import (
	"database/sql"

	"github.com/ionous/errutil"
)

func QueryAll(db *sql.DB, q string, cb func() error, dest ...interface{}) (err error) {
	if rows, e := db.Query(q); e != nil {
		err = e
	} else {
		err = ScanAll(rows, cb, dest...)
	}
	return
}

func ScanAll(rows *sql.Rows, cb func() error, dest ...interface{}) (err error) {
	for rows.Next() {
		if e := rows.Scan(dest...); e != nil {
			err = e
			break
		} else if e := cb(); e != nil {
			err = e
			break
		}
	}
	if e := rows.Err(); e != nil {
		err = errutil.Append(err, e)
	}
	rows.Close()
	return
}
