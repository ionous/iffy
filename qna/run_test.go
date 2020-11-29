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
	"github.com/ionous/iffy/test/testdb"
)

// no idea where this test should live...
// complete, manual, end to end test of factorial pattern.
func TestFullFactorial(t *testing.T) {
	db := newQnaDB(t, testdb.Memory)
	defer db.Close()

	//import factorialStory, assemble and run.
	if cnt, e := testAll(t.Name(), db, debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else if cnt != 1 {
		t.Fatal("expected one test", cnt)
	} else {
		t.Log("ok", cnt)
	}
}

func testAll(inFile string, db *sql.DB, m reader.Map) (ret int, err error) {
	var ds assembly.Dilemmas
	if e := tables.CreateAll(db); e != nil {
		err = errutil.New("couldn't create tables", e)
	} else if e := story.ImportStory(inFile, db, m); e != nil {
		err = errutil.New("couldn't import story", e)
	} else if e := assembly.AssembleStory(db, "kinds", ds.Add); e != nil {
		err = errutil.New("couldnt assemble story", e, ds.Err())
	} else if len(ds) > 0 {
		err = errutil.New("issues assembling", ds.Err())
	} else {
		ret, err = CheckAll(db, "")
	}
	return
}
