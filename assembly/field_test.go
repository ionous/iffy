package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
)

// write some primitives
func writeFields(rec *ephemera.Recorder, kfps []kfp) (err error) {
	for _, p := range kfps {
		kind := rec.Named(ephemera.NAMED_KIND, p.kind, "test")
		field := rec.Named(ephemera.NAMED_FIELD, p.field, "test")
		rec.NewPrimitive(p.fieldType, kind, field)
	}
	return
}

// name some fields that arent otherwise referenced
func writeMissing(rec *ephemera.Recorder, missing []string) (err error) {
	for _, m := range missing {
		rec.Named(ephemera.NAMED_FIELD, m, "test")
	}
	return
}

func matchProperties(db *sql.DB, want []kfp) (err error) {
	var curr kfp
	var have []kfp
	if e := dbutil.QueryAll(db,
		`select kind,field,type from mdl_field order by kind, field, type`,
		func() (err error) {
			have = append(have, curr)
			return
		}, &curr.kind, &curr.field, &curr.fieldType); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch", "have:", pretty.Sprint(have), "want:", pretty.Sprint(want))
	}
	return
}

func TestFields(t *testing.T) {
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
		} else if e := writeFields(t.rec, []kfp{
			{"P", "a", ephemera.PRIM_TEXT},
			{"Q", "b", ephemera.PRIM_TEXT},
			{"T", "c", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineFields(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchProperties(t.db, []kfp{
			{"P", "a", ephemera.PRIM_TEXT},
			{"Q", "b", ephemera.PRIM_TEXT},
			{"T", "c", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

func TestFieldLca(t *testing.T) {
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
		} else if e := writeFields(t.rec, []kfp{
			{"P", "a", ephemera.PRIM_TEXT},
			{"Q", "a", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineFields(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchProperties(t.db, []kfp{
			{"T", "a", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestFieldTypeMismatch verifies that ephemera with conflicting primitive types generates an error
// ex. T.a:text, T.a:digi
func TestFieldTypeMismatch(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := fakeHierarchy(t.modeler, []pair{
			{"T", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeFields(t.rec, []kfp{
			{"T", "a", ephemera.PRIM_TEXT},
			{"T", "a", ephemera.PRIM_DIGI},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineFields(t.modeler, t.db); e != nil {
			t.Log("okay:", e)
		} else {
			t.Fatal("expected error")
		}
	}
}
