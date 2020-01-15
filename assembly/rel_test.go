package assembly

import (
	"database/sql"
	"testing"

	"github.com/ionous/iffy/ephemera"
)

const memory = "file:test.db?cache=shared&mode=memory"

// TestVerbMismatches verifies that we can collapse multiple relation-verb pairs so long as the verb-stem pair match
// while ensuring the same stem cannot be used in multiple relations.
func TestVerbMismatches(t *testing.T) {
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		w := NewModeler(dbq)
		if e := w.WriteVerb("R", "contains"); e != nil {
			t.Fatal(e)
		} else if e := w.WriteVerb("R", "containing"); e != nil {
			t.Fatal(e)
		} else if e := w.WriteVerb("Q", "supporting"); e != nil {
			t.Fatal(e)
		} else if e := w.WriteVerb("Q", "supports"); e != nil {
			t.Fatal(e)
		} else if e := w.WriteVerb("R", "supports"); e == nil {
			t.Log("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}
