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

// TestDefaultInvalidType
func TestDefaultInvalidType(t *testing.T) {
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
				{"T", "A", ephemera.PRIM_ASPECT},
			}); e != nil {
				err = e
			} else if e := fakeTraits(t.modeler, []pair{
				{"A", "a"},
			}); e != nil {
				err = e
			} else if e := addDefaults(t.rec, vals); e != nil {
				err = e
			} else if e := DetermineDefaults(t.modeler, t.db); e == nil {
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
	/*
		fix? bools in sqlite are stored as int64;
		consider reworking to store the strings "true" or "false".

		https://www.sqlite.org/datatype3.html#boolean_datatype

		if e := testInvalid(t, []triplet{
			{"T", "d", true},
		}); e != nil {
			t.Fatal(e)
		} else */
	if e := testInvalid(t, []triplet{
		{"T", "t", false},
	}); e != nil {
		t.Fatal(e)
	}
	// try to set text and bools to aspects
	if e := testInvalid(t, []triplet{
		{"T", "A", 1.2},
	}); e != nil {
		t.Fatal(e)
	} else if e := testInvalid(t, []triplet{
		{"T", "A", true},
	}); e != nil {
		t.Fatal(e)
	} else if e := testInvalid(t, []triplet{
		{"T", "A", "x"},
	}); e != nil {
		t.Fatal(e)
	}
	// try to set text and digits to aspects
	if e := testInvalid(t, []triplet{
		{"T", "a", 1.2},
	}); e != nil {
		t.Fatal(e)
	} else if e := testInvalid(t, []triplet{
		{"T", "a", "boop"},
	}); e != nil {
		t.Fatal(e)
	}
}
