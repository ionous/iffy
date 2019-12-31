package ephemera

import (
	"database/sql"
	"log"
	"strings"
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

func (q *DbQueue) Prep(which string, keys ...string) {
	if _, ok := q.statements[which]; ok {
		log.Fatal(which, "already exists")
	} else {
		//"insert into foo(id, name) values(?, ?)"
		cols := "(" + strings.Join(keys, ", ") + ")"
		vals := "(" + strings.Repeat("?, ", len(keys)) + ")"
		if stmt, e := q.db.Prepare("INSERT into " + which + cols + " values" + vals); e != nil {
			log.Fatal(e)
		} else {
			q.statements[which] = stmt
		}
	}
	return
}

func (q *DbQueue) Write(which string, args ...interface{}) (ret Queued) {
	if stmt, ok := q.statements[which]; !ok {
		log.Fatal(which, " doesn't exist")
	} else if res, e := stmt.Exec(args...); e != nil {
		log.Fatal(e)
	} else if id, e := res.LastInsertId(); e != nil {
		log.Fatal(e)
	} else {
		ret = Queued{id}
	}
	return
}

// for returning string primary keys
//dbConn.QueryRow("INSERT INTO product (title) VALUES ($1) RETURNING product_id", p.Title).Scan(&p)

// _, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
// stmt, err = db.Prepare("select name from foo where id = ?")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	var name string
// 	err = stmt.QueryRow("3").Scan(&name)
// sqlStmt := `
// 	create table foo (id integer not null primary key, name text);
// 	delete from foo;
// 	`
// 	_, err = db.Exec(sqlStmt)
// 	if err != nil {
// 		log.Printf("%q: %s\n", err, sqlStmt)
// 		return
// 	}

// stmt, err := db.Prepare("INSERT INTO projects(id, mascot, release, category) VALUES( ?, ?, ?, ? )")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

// 	for id, project := range projects {
// 		if _, err := stmt.Exec(id+1, project.mascot, project.release, "open source"); err != nil {
// 			log.Fatal(err)
// 		}
// 	}
