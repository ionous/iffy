package assembly

import (
	"strconv"
	"strings"
	"testing"

	"github.com/ionous/iffy/tables"
)

// does the custom sql function from pluralize.go function correctly.
func TestPluralizeKinds(t *testing.T) {
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
		if e := AssemblePlurals(asm.assembler); e != nil {
			t.Fatal(e)
		}
		var buf strings.Builder
		tables.WriteCsv(asm.db, &buf,
			`select * from mdl_plural order by one`, 2)
		if have, want := buf.String(), lines(
			"animal,animals",
			"fish,fish",
			"person,people",
			"thing,things",
		); have != want {
			t.Fatal(have)
		}
	}
}
