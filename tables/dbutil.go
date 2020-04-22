package tables

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
)

// Query used for QueryAll to hide differences b/t tables.Cache and sql.DB
type Query interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// QueryAll queries the db ( or statement cache ) for one or more rows.
// For each row, it writes the row to the 'dest' args and calls 'cb' for processing.
func QueryAll(db Query, q string, cb func() error, dest ...interface{}) (err error) {
	if rows, e := db.Query(q); e != nil {
		err = e
	} else {
		err = ScanAll(rows, cb, dest...)
	}
	return
}

// ScanAll writes each row to the 'dest' args and calls 'cb' for processing.
// It closes rows before returning.
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

// Insert creates a sqlite friendly insert statement.
// For example: "insert into foo(col1, col2, ...) values(?, ?, ...)"
func Insert(table string, keys ...string) string {
	vals := "?"
	if kcnt := len(keys) - 1; kcnt > 0 {
		vals += strings.Repeat(",?", kcnt)
	}
	return "INSERT into " + table +
		"(" + strings.Join(keys, ", ") + ")" +
		" values " + "(" + vals + ");"
}
