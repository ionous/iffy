package qna

import (
	"database/sql"
	"os/user"
	"path"
	"testing"

	"github.com/ionous/iffy/assembly"
)

const memory = "file:test.db?cache=shared&mode=memory"

// if path is nil, it will use a file db.
func newQnaDB(t *testing.T, where string) (ret *sql.DB) {
	var source string
	if len(where) > 0 {
		source = where
	} else if user, e := user.Current(); e != nil {
		t.Fatal(e)
	} else {
		source = path.Join(user.HomeDir, t.Name()+".db")
	}
	//
	if db, e := sql.Open(assembly.SqlCustomDriver, source); e != nil {
		t.Fatal(e)
	} else {
		ret = db
	}
	return
}
