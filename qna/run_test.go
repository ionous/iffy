package qna

import (
	"database/sql"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/ephemera/story"
	"github.com/ionous/iffy/tables"
)

// no idea where this test should live...
// complete, manual, end to end test of factorial pattern.
func TestFullFactorial(t *testing.T) {
	db := newTestDB(t, memory)
	defer db.Close()

	//import factorialStory, assemble and run.
	if e := testAll(t.Name(), db, debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else {
		t.Log("ok")
	}
}

func testAll(inFile string, db *sql.DB, m reader.Map) (err error) {
	var ds assembly.Dilemmas
	if e := tables.CreateAll(db); e != nil {
		err = errutil.New("couldn't create tables", e)
	} else if e := story.ImportStory(inFile, db, m); e != nil {
		err = errutil.New("couldn't import story", e)
	} else if e := assembly.AssembleStory(db, "things", ds.Add); e != nil {
		err = errutil.New("couldnt assemble story", e)
	} else if len(ds) > 0 {
		err = errutil.New("issues assembling", ds.Err())
	} else if e := CheckAll(db); e != nil {
		err = e
	}
	return
}
