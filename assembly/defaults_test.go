package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/kr/pretty"
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
		} else if e := DetermineDefaults(t.modeler, t.db); e != nil {
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
			{"Q", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeTraits(t.modeler, []pair{
			{"A", "w"},
			{"A", "x"},
			{"A", "y"},
			{"B", "z"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeAspects(t.modeler, []pair{
			{"T", "A"},
			{"P", "B"},
			{"Q", "B"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addDefaults(t.rec, []triplet{
			{"T", "x", true},
			{"P", "y", true},
			{"P", "z", true},
			//
			{"Q", "A", "w"},
			{"Q", "B", "z"},
			{"Q", "w", true},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineDefaults(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchDefaults(t.db, []triplet{
			{"P", "A", "y"},
			{"P", "B", "z"},
			{"Q", "A", "w"},
			{"Q", "B", "z"},
			{"T", "A", "x"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultDuplicates to verify that duplicate default specifications are okay
func TestDefaultDuplicates(t *testing.T) {
	if t, e := newDefaultsTest(t, memory, []triplet{
		{"T", "t", "text"},
		{"T", "t", "text"},
		{"P", "t", "text"},
		//
		{"T", "d", 123},
		{"T", "d", 123},
		{"P", "d", 123},
		//
		{"T", "A", "y"},
		{"T", "y", true},
		{"P", "x", true},
		{"P", "A", "x"},
	}); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := DetermineDefaults(t.modeler, t.db); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultConflict to verify that conflicting values for the same default are not okay
func TestDefaultConflict(t *testing.T) {
	testConflict := func(t *testing.T, vals []triplet) (err error) {
		if t, e := newDefaultsTest(t, memory, vals); e != nil {
			t.Fatal(e)
		} else {
			defer t.Close()
			if e := DetermineDefaults(t.modeler, t.db); e == nil {
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

	if e := testConflict(t, []triplet{
		{"T", "A", "x"},
		{"T", "A", "z"},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []triplet{
		{"T", "x", true},
		{"T", "z", true},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []triplet{
		{"T", "A", "x"},
		{"T", "z", true},
	}); e != nil {
		t.Fatal(e)
	}
}

// TestDefaultBadValue to verify that modeling requires appropriate values for defaults based on type
func TestDefaultBadValue(t *testing.T) {
	//- for now, we only allow text and number [ text and digi ]
	// - later we could add ambiguity for conversion [ 4 -> "4" ]
	testInvalid := func(t *testing.T, vals []triplet) (err error) {
		if t, e := newDefaultsTest(t, memory, vals); e != nil {
			err = e
		} else {
			defer t.Close()
			if e := DetermineDefaults(t.modeler, t.db); e == nil {
				err = errutil.New("expected error")
			} else {
				t.Log("okay:", e)
			}
		}
		return
	}

	if e := testInvalid(t, []triplet{
		{"T", "t", 1.2},
	}); e != nil {
		t.Fatal(e)
	} else if e := testInvalid(t, []triplet{
		{"T", "d", "1.2"},
	}); e != nil {
		t.Fatal(e)
	}
	// try to set trait like values

	if e := testInvalid(t, []triplet{
		{"T", "t", false},
	}); e != nil {
		t.Fatal(e)
	}

	/*
	   fix? somehow? bools in sqlite are stored as int64;
	   could switch to text ( "true", "false" ) perhaps and add some check/query
	   during determination
	   if e := testInvalid(t, []triplet{
	       {"T", "d", true},
	   }); e != nil {
	       t.Fatal(e)
	   }
	*/

	/* fix? aspects are set by matching traits
	1.2 is not a trait, so it's skipped.
	this might get handled by a "missing" check,
	or possibly by changing the determination query.

	if e := testInvalid(t, []triplet{
		{"T", "A", 1.2},
	}); e != nil {
		t.Fatal(e)
	}
	*/
}

// match generated model defaults
func matchDefaults(db *sql.DB, want []triplet) (err error) {
	var curr triplet
	var have []triplet
	if e := dbutil.QueryAll(db,
		`select kind, field, value 
			from mdl_default
			order by kind, field, value`,
		func() (err error) {
			have = append(have, curr)
			return
		},
		&curr.target, &curr.prop, &curr.value); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch",
			"have:", pretty.Sprint(have),
			"want:", pretty.Sprint(want))
	}
	return
}

// write ephemera describing some initial values
func addDefaults(rec *ephemera.Recorder, defaults []triplet) (err error) {
	for _, el := range defaults {
		namedKind := rec.Named(ephemera.NAMED_KIND, el.target, "test")
		namedField := rec.Named(ephemera.NAMED_PROPERTY, el.prop, "test")
		rec.NewDefault(namedKind, namedField, el.value)
	}
	return
}

func newDefaultsTest(t *testing.T, path string, defaults []triplet) (ret *assemblyTest, err error) {
	if t, e := newAssemblyTest(t, path); e != nil {
		err = e
	} else {
		if e := fakeHierarchy(t.modeler, []pair{
			{"T", ""},
			{"P", "T"},
		}); e != nil {
			err = e
		} else if e := fakeFields(t.modeler, []kfp{
			{"T", "d", ephemera.PRIM_DIGI},
			{"T", "t", ephemera.PRIM_TEXT},
			{"T", "A", ephemera.PRIM_ASPECT},
		}); e != nil {
			err = e
		} else if e := fakeTraits(t.modeler, []pair{
			{"A", "x"}, {"A", "y"}, {"A", "z"},
		}); e != nil {
			err = e
		} else if e := addDefaults(t.rec, defaults); e != nil {
			err = e
		}
		if err != nil {
			t.Close()
		} else {
			ret = t
		}
	}
	return
}
