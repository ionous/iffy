package assembly

import (
	"strconv"
	"testing"

	"github.com/ionous/iffy/tables"
)

// does the custom sql function from pluralize.go function correctly.
func TestSqlPluralize(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		for i, n := range []string{
			"thing",
			"animal",
			"person",
			"fish",
		} {
			asm.rec.NewName(n, tables.NAMED_KIND, strconv.Itoa(i))
		}
		var plurals []string
		var plural string
		if e := tables.QueryAll(asm.db,
			`select pluralize(name) 
			from eph_named
			where category='singular_kind'
			order by name`,
			func() error {
				plurals = append(plurals, plural)
				return nil
			},
			&plural); e != nil {
			t.Fatal(e)
		}
		if want := lines(
			"animals",
			"fish",
			"people",
			"things",
		); want != lines(plurals...) {
			t.Fatal(plurals)
		}
		// asm.db.
	}
}
