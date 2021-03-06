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

func Must(db *sql.DB, q string, args ...interface{}) (ret int64) {
	if res, e := db.Exec(q, args...); e != nil {
		panic(e)
	} else if id, e := res.LastInsertId(); e != nil {
		panic(e)
	} else {
		ret = id
	}
	return
}

func RowsAffected(res sql.Result) (ret int) {
	if cnt, e := res.RowsAffected(); e != nil {
		ret = -1
	} else {
		ret = int(cnt)
	}
	return
}

// QueryAll queries the db ( or statement cache ) for one or more rows.
// For each row, it writes the row to the 'dest' args and calls 'cb' for processing.
func QueryAll(db Query, q string, cb func() error, dest ...interface{}) (err error) {
	if rows, e := db.Query(q); e != nil {
		err = errutil.New("QueryAll error:", e)
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
			err = errutil.New("ScanAll error:", e)
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
	return InsertWith(table, "", keys...)
}

func InsertWith(table string, rest string, keys ...string) string {
	vals := "?"
	if kcnt := len(keys) - 1; kcnt > 0 {
		vals += strings.Repeat(",?", kcnt)
	}
	return "INSERT into " + table +
		"(" + strings.Join(keys, ", ") + ")" +
		" values " + "(" + vals + ")" + rest + ";"
}
