package tables

import "database/sql"

type Cache struct {
	db    *sql.DB
	cache map[string]*sql.Stmt
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
