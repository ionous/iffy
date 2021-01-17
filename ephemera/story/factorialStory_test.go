package story_test

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/story"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

// read the factorial test story
func TestFactorialStory(t *testing.T) {
	db := newImportDB(t, testdb.Memory)
	defer db.Close()
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create tables", e)
	} else if e := debug.FactorialStory.ImportStory(story.NewImporter(t.Name(), db)); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select count() from eph_check", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_rule", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_prog", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_pattern", 1)
		tables.WriteCsv(db, &buf, "select distinct name, category from eph_named where category != 'scene' order by name, category", 2)
		if have, want := buf.String(), lines(
			"1", // eph_check -- 1 unit test
			"2", // eph_rule -- 2 rules
			"3", // e@h_prog -- 1 test program, 2 rules
			"4", // eph_pattern specifies types - (1 pattern, 1 parameter) * (1 decl, 1 call)
			// eph_named
			"execute,type",      // we evaluated the return
			"factorial,pattern", // name of the pattern
			"factorial,test",    // name of the test
			"num,parameter",     // we declared the param
			"num,return",        // we referenced the return
			"number_eval,type",  // we evaluated the var
		); have != want {
			t.Fatal(have)
		}
	}
}
