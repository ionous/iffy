package tables

import (
	"database/sql"

	"github.com/ionous/errutil"
)

// Cache mimics the sql.Stmt api, creating the Stmt objects on demand.
type Cache struct {
	db    *sql.DB
	cache map[string]*sql.Stmt
}

type prepError struct {
	err error
}

// RowScanner because sql.Row.Scan doesnt have the sql.Scanner.Scan interface.
type RowScanner interface {
	Scan(...interface{}) error
}

// implements RowScanner
func (e *prepError) Scan(...interface{}) error {
	return e.err

}

func NewCache(db *sql.DB) *Cache {
	return &Cache{db, make(map[string]*sql.Stmt)}
}

func (c *Cache) DB() *sql.DB {
	return c.db
}

func (c *Cache) Close() {
	for _, v := range c.cache {
		v.Close()
	}
	c.cache = make(map[string]*sql.Stmt)
}

func (c *Cache) Must(q string, args ...interface{}) (ret int64) {
	if id, e := c.Exec(q, args...); e != nil {
		panic(e)
	} else {
		ret = id
	}
	return
}

func (c *Cache) Exec(q string, args ...interface{}) (ret int64, err error) {
	if stmt, e := c.prep(q); e != nil {
		err = e
	} else if res, e := stmt.Exec(args...); e != nil {
		err = errutil.New("Exec error:", e)
	} else if id, e := res.LastInsertId(); e != nil {
		err = errutil.New("LastInsertId error:", e)
	} else {
		ret = id
	}
	return
}

func (c *Cache) Query(q string, args ...interface{}) (ret *sql.Rows, err error) {
	if stmt, e := c.prep(q); e != nil {
		err = e
	} else if rows, e := stmt.Query(args...); e != nil {
		err = errutil.New("Query error", e)
	} else {
		ret = rows
	}
	return
}

// QueryRow mimics db.QueryRow but returns Scanner instead of Row
// so that we can defer any errors encountered while preparing the cached statement.
func (c *Cache) QueryRow(q string, args ...interface{}) (ret RowScanner) {
	if stmt, e := c.prep(q); e != nil {
		ret = &prepError{e}
	} else {
		ret = stmt.QueryRow(args...)
	}
	return
}

func (c *Cache) prep(q string) (ret *sql.Stmt, err error) {
	if c.db == nil {
		err = errutil.New("cache not initialized")
	} else if stmt := c.cache[q]; stmt != nil {
		ret = stmt
	} else if stmt, e := c.db.Prepare(q); e != nil {
		err = errutil.New("Prepare error:", e)
	} else {
		c.cache[q] = stmt
		ret = stmt
	}
	return
}
