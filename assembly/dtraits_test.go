package assembly

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestDefaultTraitAssignment to verify default traits can be assigned to kinds.
func TestDefaultTraitAssignment(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := fakeHierarchy(t.modeler, []pair{
			{"T", ""},
			{"P", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeTraits(t.modeler, []pair{
			{"A", "x"},
			{"A", "y"},
			{"B", "z"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeAspects(t.modeler, []pair{
			{"T", "A"},
			{"P", "B"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addDefaults(t.rec, []triplet{
			{"T", "x", true},
			{"P", "y", true},
			{"P", "z", true},
		}); e != nil {
			t.Fatal(e)
		} else if e := determineDefaultTraits(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchDefaults(t.db, []triplet{
			{"P", "A", "y"},
			{"P", "B", "z"},
			{"T", "A", "x"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}
