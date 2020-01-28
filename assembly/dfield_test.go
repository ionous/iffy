package assembly

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	_ "github.com/mattn/go-sqlite3"
)

// TestDefaultFieldAssigment to verify default values can be assigned to kinds.
func TestDefaultFieldAssigment(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := fakeHierarchy(t.modeler, []pair{
			{"T", ""},
			{"P", "T"},
			{"D", "T"},
			{"C", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeFields(t.modeler, []kfp{
			{"T", "d", ephemera.PRIM_DIGI},
			{"T", "t", ephemera.PRIM_TEXT},
			{"T", "t2", ephemera.PRIM_TEXT},
			{"P", "x", ephemera.PRIM_TEXT},
			{"D", "x", ephemera.PRIM_TEXT},
			{"C", "c", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := addDefaults(t.rec, []triplet{
			{"T", "t", "some text"},
			{"P", "t", "override text"},
			{"P", "t2", "other text"},
			{"P", "x", "x in p"},
			{"D", "x", "x in d"},
			{"C", "c", "c text"},
			{"C", "d", 123},
		}); e != nil {
			t.Fatal(e)
		} else if e := determineDefaultFields(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchDefaults(t.db, []triplet{
			{"C", "c", "c text"},
			{"C", "d", int64(123)}, // re: int64 -- default scanner uses https://golang.org/pkg/database/sql/#Scanner
			{"D", "x", "x in d"},
			{"P", "t", "override text"},
			{"P", "t2", "other text"},
			{"P", "x", "x in p"},
			{"T", "t", "some text"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultFieldDuplicate to verify that duplicate values are okay
func TestDefaultFieldDuplicate(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := fakeHierarchy(t.modeler, []pair{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeFields(t.modeler, []kfp{
			{"T", "d", ephemera.PRIM_DIGI},
			{"T", "t", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := addDefaults(t.rec, []triplet{
			{"T", "t", "text"},
			{"T", "t", "text"},
			{"C", "t", "text"},
			//
			{"T", "d", 123},
			{"T", "d", 123},
			{"C", "d", 123},
		}); e != nil {
			t.Fatal(e)
		} else if e := determineDefaultFields(t.modeler, t.db); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultFieldsConflict to verify that conflicting values are not okay
func TestDefaultFieldsConflict(t *testing.T) {
	testConflict := func(t *testing.T, vals []triplet) (err error) {
		if t, e := newAssemblyTest(t, memory); e != nil {
			err = e
		} else {
			defer t.Close()
			//
			if e := fakeHierarchy(t.modeler, []pair{
				{"T", ""},
				{"P", "T"},
			}); e != nil {
				t.Fatal(e)
			} else if e := fakeFields(t.modeler, []kfp{
				{"T", "d", ephemera.PRIM_DIGI},
				{"T", "t", ephemera.PRIM_TEXT},
			}); e != nil {
				err = e
			} else if e := addDefaults(t.rec, vals); e != nil {
				err = e
			} else if e := determineDefaultFields(t.modeler, t.db); e == nil {
				err = errutil.New("expected error")
			} else {
				t.Log("okay:", e)
			}
		}
		return
	}
	if e := testConflict(t, []triplet{
		{"T", "t", "a"},
		{"T", "t", "b"},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []triplet{
		{"T", "d", 1},
		{"T", "d", 2},
	}); e != nil {
		t.Fatal(e)
	}
}
