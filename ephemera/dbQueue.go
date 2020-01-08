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
		desc[i] = c.Name + "  " + c.Type
	}
	create := "drop table if exists " + which + "; create table " + which + "(" + strings.Join(desc, ", ") + ");"
	if _, e := q.db.Exec(create); e != nil {
		err = errutil.New(e, "for", create)
	}
	return
}

func (q *DbQueue) Prep(which string, cols ...Col) {
	if _, ok := q.statements[which]; ok {
		log.Fatalln("prep", which, "already exists")
	} else if e := q.Create(which, cols); e != nil {
		log.Fatalln("prep", e)
	} else {
		//
		//"insert into foo(id, name) values(?, ?)"
		vals := "?"
		if kcnt := len(cols) - 1; kcnt > 0 {
			vals += strings.Repeat(",?", kcnt)
		}
		keys := make([]string, len(cols))
		for i, c := range cols {
			keys[i] = c.Name
		}
		query := "INSERT into " + which +
			"(" + strings.Join(keys, ", ") + ")" +
			" values " + "(" + vals + ");"
		if stmt, e := q.db.Prepare(query); e != nil {
			log.Fatalln("prep insert", query, e)
		} else {
			q.statements[which] = stmt
		}
	}
	return
}

func (q *DbQueue) Write(which string, args ...interface{}) (ret Queued) {
	if stmt, ok := q.statements[which]; !ok {
		log.Fatalln(which, "doesn't exist")
	} else if res, e := stmt.Exec(args...); e != nil {
		log.Fatalln("write", e)
	} else if id, e := res.LastInsertId(); e != nil {
		log.Fatalln("last id", e)
	} else {
		ret = Queued{id}
	}
	return
}
