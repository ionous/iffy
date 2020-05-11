package story

import (
	"database/sql"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera/debug"
)

// test calling a pattern
// note: the pattern is undefined.
func TestFactorialStory(t *testing.T) {
	errutil.Panic = true
	const memory = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal("db open", e)
	} else {
		defer db.Close()
		if e := ImportStory(t.Name(), debug.FactorialStory, db); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}
