package assembly

import (
	"database/sql"
	"os/user"
	"path"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
)

type Output struct {
}

func (out *Output) Conflict(e error) {
}
func (out *Output) Ambiguity(e error) {
}

func getPath() (ret string, err error) {
	if user, e := user.Current(); e != nil {
		err = e
	} else {
		ret = path.Join(user.HomeDir, "test.db")
	}
	return
}

func MissingKinds(db *sql.DB, cb func(kind string)) (err error) {
	if kinds, e := db.Query(
		`select distinct named.name as kind 
			 from named
			 where kind not in (
			 	select kind from ancestry
			 ) and named.category = 'kind'`); e != nil {
		err = e
	} else {
		for kinds.Next() {
			var k string
			if e := kinds.Scan(&k); e != nil {
				err = e
				break
			} else {
				cb(k)
			}
		}
		if e := kinds.Err(); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func NewWriter(q ephemera.Queue) *Writer {
	q.Prep("Ancestry",
		ephemera.Col{"kind", "text"},
		ephemera.Col{"path", "text"})
	return &Writer{q}
}

type Writer struct {
	q ephemera.Queue
}

// write kind and comma separated ancestors
func (w *Writer) WriteAncestor(kind, path string) {
	// fix error
	w.q.Write("Ancestry", kind, path)
}
