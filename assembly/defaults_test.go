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

func TestDefaultsAssigment(t *testing.T) {
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
			{"T", "t2", ephemera.PRIM_TEXT},
			{"P", "p", ephemera.PRIM_TEXT},
			{"C", "c", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeDefaults(t.rec, []prop{
			{"T", "t", "some text"},
			{"P", "t2", "other text"},
			{"C", "c", "c text"},
			{"C", "d", 123},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineDefaults(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchDefaults(t.db, []prop{
			{"C", "c", "c text"},
			// default scanner uses https://golang.org/pkg/database/sql/#Scanner
			{"T", "d", int64(123)}, //
			{"T", "t", "some text"},
			{"T", "t2", "other text"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultsUnknownField missing properties ( kind, field pair doesn't exist in model )
func TestDefaultsUnknownField(t *testing.T) {
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
			{"T", "t2", ephemera.PRIM_TEXT},
			{"P", "p", ephemera.PRIM_TEXT},
			{"C", "c", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeDefaults(t.rec, []prop{
			{"T", "t", "some text"},
			{"P", "t2", "other text"},
			{"C", "c", "c text"},
			{"P", "c", "invalid"}, // this pair doesnt exist
			{"T", "p", "invalid"}, // this pair doesnt exist
			{"C", "d", 123},
		}); e != nil {
			t.Fatal(e)
		} else {
			var got []pair
			if e := MissingDefaults(t.db, func(k, f string) (err error) {
				got = append(got, pair{k, f})
				return
			}); e != nil {
				t.Fatal(e)
			} else if !reflect.DeepEqual(got, []pair{
				{"P", "c"},
				{"T", "p"},
			}) {
				t.Fatal("mismatched", got)
			}
		}
	}
}

// TestDefaultsValuesDuplicate to verify that duplicate values are okay
func TestDefaultsValuesDuplicate(t *testing.T) {
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
		} else if e := writeDefaults(t.rec, []prop{
			{"T", "t", "text"},
			{"P", "t", "text"},
			{"C", "t", "text"},
			//
			{"T", "d", 123},
			{"P", "d", 123},
			{"C", "d", 123},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineDefaults(t.modeler, t.db); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultsValuesConflict to verify that conflicting values are not okay
func TestDefaultsValuesConflict(t *testing.T) {
	testConflict := func(t *testing.T, vals []prop) (err error) {
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
			} else if e := writeDefaults(t.rec, vals); e != nil {
				err = e
			} else if e := DetermineDefaults(t.modeler, t.db); e == nil {
				err = errutil.New("expected error")
			} else {
				t.Log("okay:", e)
			}
		}
		return
	}

	if e := testConflict(t, []prop{
		{"T", "t", "a"},
		{"T", "t", "b"},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []prop{
		{"T", "t", "a"},
		{"P", "t", "b"},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []prop{
		{"T", "d", 1},
		{"T", "d", 2},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []prop{
		{"T", "d", 1},
		{"P", "d", 2},
	}); e != nil {
		t.Fatal(e)
	}
}

// TestDefaultsInvalidType
func TestDefaultsInvalidType(t *testing.T) {
	//- for now, we only allow text and number [ text and digi ]
	// - later we could add ambiguity for conversion [ 4 -> "4" ]
	testInvalid := func(t *testing.T, vals []prop) (err error) {
		if t, e := newAssemblyTest(t, memory); e != nil {
			err = e
		} else {
			defer t.Close()
			//
			if e := fakeHierarchy(t.modeler, []pair{
				{"T", ""},
			}); e != nil {
				t.Fatal(e)
			} else if e := fakeFields(t.modeler, []kfp{
				{"T", "d", ephemera.PRIM_DIGI},
				{"T", "t", ephemera.PRIM_TEXT},
			}); e != nil {
				err = e
			} else if e := writeDefaults(t.rec, vals); e != nil {
				err = e
			} else if e := DetermineDefaults(t.modeler, t.db); e == nil {
				err = errutil.New("expected error")
			} else {
				t.Log("okay:", e)
			}
		}
		return
	}
	if e := testInvalid(t, []prop{
		{"T", "t", 1.2},
	}); e != nil {
		t.Fatal(e)
	} else if e := testInvalid(t, []prop{
		{"T", "d", "1.2"},
	}); e != nil {
		t.Fatal(e)
	}
}

// match generated model defaults
func matchDefaults(db *sql.DB, want []prop) (err error) {
	var curr prop
	var have []prop
	if e := dbutil.QueryAll(db,
		`select kind, field, value 
			from mdl_default
			order by kind, field, value`,
		func() (err error) {
			have = append(have, curr)
			return
		},
		&curr.owner, &curr.prop, &curr.value); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch",
			"have:", pretty.Sprint(have),
			"want:", pretty.Sprint(want))
	}
	return
}

// write ephemera describing some initial values
func writeDefaults(rec *ephemera.Recorder, defaults []prop) (err error) {
	for _, el := range defaults {
		namedKind := rec.Named(ephemera.NAMED_KIND, el.owner, "test")
		namedField := rec.Named(ephemera.NAMED_FIELD, el.prop, "test")
		rec.NewDefault(namedKind, namedField, el.value)
	}
	return
}
