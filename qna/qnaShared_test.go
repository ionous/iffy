package qna

import (
	"database/sql"
	"testing"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/test/testdb"
)

// if path is nil, it will use a file db.
func newQnaDB(t *testing.T, path string) (ret *sql.DB) {
	var source string
	if len(path) > 0 {
		source = path
	} else if p, e := testdb.PathFromName(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		source = p
	}
	// assembly needed for TestFullFactorial
	if db, e := sql.Open(assembly.SqlCustomDriver, source); e != nil {
		t.Fatal(e)
	} else {
		ret = db
	}
	return
}
