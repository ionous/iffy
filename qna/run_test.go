package qna

import (
	"testing"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/ephemera/story"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

// no idea where this test should live...
// tests the execution of an imported story;
// doesnt test the *reading* of the story
func TestFullFactorial(t *testing.T) {
	db := newQnaDB(t, testdb.Memory)
	defer db.Close()

	//import factorialStory, assemble and run.
	var ds reader.Dilemmas
	if e := tables.CreateAll(db); e != nil {
		t.Fatal("couldn't create tables", e)
	} else if e := debug.FactorialStory.ImportStory(story.NewImporter(t.Name(), db)); e != nil {
		t.Fatal("couldn't import story", e)
	} else if e := assembly.AssembleStory(db, "kinds", ds.Add); e != nil {
		t.Fatal("couldnt assemble story", e, ds.Err())
	} else if len(ds) > 0 {
		t.Fatal("issues assembling", ds.Err())
	} else if cnt, e := CheckAll(db, ""); e != nil {
		t.Fatal(e)
	} else if cnt != 1 {
		t.Fatal("expected one test", cnt)
	} else {
		t.Log("ok", cnt)
	}
}
