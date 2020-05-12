package qna

import (
	"database/sql"
	"os/user"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const memory = "file:test.db?cache=shared&mode=memory"

// if path is nil, it will use a file db.
func newTestDB(t *testing.T, where string) (ret *sql.DB) {
	var source string
	if len(where) > 0 {
		source = where
	} else if user, e := user.Current(); e != nil {
		t.Fatal(e)
	} else {
		source = path.Join(user.HomeDir, t.Name()+".db")
	}
	//
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		ret = db
	}
	return
}
