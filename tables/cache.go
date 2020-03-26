package tables

import "database/sql"

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
		err = e
	} else if id, e := res.LastInsertId(); e != nil {
		err = e
	} else {
		ret = id
	}
	return
}

func (c *Cache) Query(q string, args ...interface{}) (ret *sql.Rows, err error) {
	if stmt, e := c.prep(q); e != nil {
		err = e
	} else {
		ret, err = stmt.Query(args...)
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
	if stmt := c.cache[q]; stmt != nil {
		ret = stmt
	} else if stmt, e := c.db.Prepare(q); e != nil {
		err = e
	} else {
		c.cache[q] = stmt
		ret = stmt
	}
	return
}
