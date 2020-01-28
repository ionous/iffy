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

// TestDefaultFieldUnknownField missing properties ( kind, field pair doesn't exist in model )
func xTestDefaultFieldUnknownField(t *testing.T) {
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
		} else if e := addDefaults(t.rec, []triplet{
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

// TestDefaultFieldValuesDuplicate to verify that duplicate values are okay
func TestDefaultFieldValuesDuplicate(t *testing.T) {
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

// TestDefaultFieldValuesConflict to verify that conflicting values are not okay
func TestDefaultFieldValuesConflict(t *testing.T) {
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

// TestDefaultFieldInvalidType
func TestDefaultFieldInvalidType(t *testing.T) {
	//- for now, we only allow text and number [ text and digi ]
	// - later we could add ambiguity for conversion [ 4 -> "4" ]
	testInvalid := func(t *testing.T, vals []triplet) (err error) {
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
	if e := testInvalid(t, []triplet{
		{"T", "t", 1.2},
	}); e != nil {
		t.Fatal(e)
	} else if e := testInvalid(t, []triplet{
		{"T", "d", "1.2"},
	}); e != nil {
		t.Fatal(e)
	}
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
