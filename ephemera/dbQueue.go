package ephemera

import (
	"database/sql"
	"log"
	"strings"

	"github.com/ionous/errutil"
)

type DbQueue struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewDBQueue(db *sql.DB) *DbQueue {
	return &DbQueue{
		db:         db,
		statements: make(map[string]*sql.Stmt),
	}
}

const AutoKey = "ROWID"

func (q *DbQueue) Create(which string, cols []Col) (err error) {
	desc := make([]string, len(cols))
	for i, c := range cols {
		desc[i] = c.Name + "  " + c.Type + " " + c.Check
	}
	create := "drop table if exists " + which + "; create table " + which + "(" + strings.Join(desc, ", ") + ");"
	if _, e := q.db.Exec(create); e != nil {
		err = errutil.New(e, "for", create)
	}
	return
}

// Prep creates a new table (which) and prepares a insert statement for it.
func (q *DbQueue) Prep(which string, cols ...Col) {
	// build a string like: "insert into foo(col1, col2, ...) values(?, ?, ...)"
	keys := NamesOf(cols)
	vals := "?"
	if kcnt := len(keys) - 1; kcnt > 0 {
		vals += strings.Repeat(",?", kcnt)
	}
	q.PrepStatement(which, "INSERT into "+which+
		"("+strings.Join(keys, ", ")+")"+
		" values "+"("+vals+");", cols)
}

func (q *DbQueue) PrepStatement(which, query string, cols []Col) {
	if _, ok := q.statements[which]; ok {
		log.Fatalln("prep", which, "already exists")
	} else if e := q.Create(which, cols); e != nil {
		log.Fatalln("prep", e)
	} else if stmt, e := q.db.Prepare(query); e != nil {
		log.Fatalln("prep insert", query, e)
	} else {
		q.statements[which] = stmt
	}
}

func (q *DbQueue) Write(which string, args ...interface{}) (ret Queued, err error) {
	if stmt, ok := q.statements[which]; !ok {
		err = errutil.New(which, "doesn't exist")
	} else if res, e := stmt.Exec(args...); e != nil {
		err = errutil.New("write", e)
	} else if id, e := res.LastInsertId(); e != nil {
		err = errutil.New("last id", e)
	} else {
		ret = Queued{id}
	}
	return
}
