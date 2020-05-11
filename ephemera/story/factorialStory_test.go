package story

import (
	"strings"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/tables"
)

// test calling a pattern
// note: the pattern is undefined.
func TestFactorialStory(t *testing.T) {
	errutil.Panic = true
	db := newTestDB(t, memory)
	defer db.Close()
	if e := ImportStory(t.Name(), debug.FactorialStory, db); e != nil {
		t.Fatal("import", e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select count() from eph_check", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_rule", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_prog", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_pattern", 1)
		tables.WriteCsv(db, &buf, "select distinct name, category from eph_named order by name", 2)
		if have, want := buf.String(), lines(
			"1", // eph_check -- 1 unit test
			"2", // eph_rule -- 2 rules
			"3", // e@h_prog -- 1 test program, 2 rules
			"5", // eph_pattern specifies types - 1 pattern, 1 parameter, 1 call, 2 rules
			// eph_named
			"assignment,type",         //--> FIX? its possible this should be number_eval
			"factorial,test",          // name of the test
			"factorial,determine_num", // we called the named pattern
			"factorial,pattern_name",  // we declared the named pattern
			"num,variable_name",       // we referenced the variable named
			"number_eval,type",        // we evaluated the pattern
		); have != want {
			t.Fatal(have)
		}
	}
}
