package story

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// test calling a pattern
// note: the pattern is undefined.
func TestDetermineNum(t *testing.T) {
	expect := pattern.DetermineNum{
		"factorial", pattern.NewNamedParams(
			"num", &core.FromNum{
				&core.Number{3},
			})}
	k, db := newTestDecoder(t, memory)
	defer db.Close()
	if rule, e := imp_determine_num(k, debug.FactorialDetermineNum); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(rule, &expect); len(diff) != 0 {
		t.Fatal(diff)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select count() from eph_prog", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_rule", 1)
		tables.WriteCsv(db, &buf, "select * from eph_pattern", 4)
		tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		if have, want := buf.String(), lines(
			// eph_prog count
			// no programs b/c no container for the call into determine.
			"0",
			// eph_rule count
			// no rules b/c the pattern is called but not implemented.
			"0",
			// eph_pattern
			"2,3,4,0", // from NewPatternRef -> "determine num" takes a parameter that is from a number eval
			"2,2,5,0", // from NewPatternRef -> "determine num" indicates factorial returns a number eval
			//
			"factorial,pattern_name", // 1.
			"num,variable_name",      // 2.
			"number_eval,type",       // 3.
			"number_eval,type",       // 4.
		); have != want {
			t.Fatal(have)
		}
	}
}
