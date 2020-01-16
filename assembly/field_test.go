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

func TestLca(t *testing.T) {
	match := func(a, b, c []string) bool {
		_, chain := findOverlap(a, b)
		return reflect.DeepEqual(chain, c)
	}
	if !match([]string{"A"}, []string{"A"}, []string{"A"}) {
		t.Fatal("expected lowest common ancestor A")
	} else if !match([]string{"A"}, []string{"B", "A"}, []string{"A"}) {
		t.Fatal("expected lowest common ancestor A")
	} else if !match([]string{"B", "A"}, []string{"B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor A")
	} else if !match([]string{"D", "C", "B", "A"}, []string{"B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor B")
	} else if !match([]string{"B", "A"}, []string{"D", "C", "B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor B")
	} else if !match([]string{"E", "F", "B", "A"}, []string{"D", "C", "B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor B")
	} else if !match([]string{"D", "C", "B", "A"}, []string{"E", "F", "B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor B")
	} else if !match([]string{"D", "E", "F"}, []string{"C", "B", "A"}, nil) {
		t.Fatal("expected no lowest common ancestor")
	}
}

type kfp struct{ kind, field, fieldType string }
type pair struct{ key, value string }

// create some fake hierarchy
func writeHierarchy(w *Modeler, kinds []pair) (err error) {
	for _, p := range kinds {
		if e := w.WriteAncestor(p.key, p.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func writeFields(db *sql.DB, kinds []pair, kfps []kfp, missing ...string) (err error) {
	dbq := ephemera.NewDBQueue(db)
	r := ephemera.NewRecorder("ancestorTest", dbq)
	w := NewModeler(dbq)
	if e := writeHierarchy(w, kinds); e != nil {
		err = e
	} else {
		// write some primitives
		for _, p := range kfps {
			kind := r.Named(ephemera.NAMED_KIND, p.kind, "test")
			field := r.Named(ephemera.NAMED_FIELD, p.field, "test")
			r.NewPrimitive(p.fieldType, kind, field)
		}
		// name some fields that arent otherwise referenced
		for _, m := range missing {
			r.Named(ephemera.NAMED_FIELD, m, "test")
		}
		if e := DetermineFields(w, db); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func matchProperties(db *sql.DB, want []kfp) (err error) {
	var curr kfp
	var have []kfp
	if e := dbutil.QueryAll(db,
		`select kind,field,type from mdl_property order by kind, field, type`,
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
	const source = memory
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeFields(db,
			[]pair{
				{"T", ""},
				{"P", "T"},
				{"Q", "T"},
			},
			[]kfp{
				{"P", "a", ephemera.PRIM_TEXT},
				{"Q", "b", ephemera.PRIM_TEXT},
				{"T", "c", ephemera.PRIM_TEXT},
			}); e != nil {
			t.Fatal(e)
		} else if e := matchProperties(db, []kfp{
			{"P", "a", ephemera.PRIM_TEXT},
			{"Q", "b", ephemera.PRIM_TEXT},
			{"T", "c", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

func TestFieldLca(t *testing.T) {
	const source = memory
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeFields(db,
			[]pair{
				{"T", ""},
				{"P", "T"},
				{"Q", "T"},
			},
			[]kfp{
				{"P", "a", ephemera.PRIM_TEXT},
				{"Q", "a", ephemera.PRIM_TEXT},
			}); e != nil {
			t.Fatal(e)
		} else if e := matchProperties(db, []kfp{
			{"T", "a", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestFieldTypeMismatch verifies that ephemera with conflicting primitive types generates an error
// ex. T.a:text, T.a:digi
func TestFieldTypeMismatch(t *testing.T) {
	const source = memory
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeFields(db,
			[]pair{{"T", ""}},
			[]kfp{
				{"T", "a", ephemera.PRIM_TEXT},
				{"T", "a", ephemera.PRIM_DIGI},
			}); e != nil {
			t.Log("okay:", e)
		} else {
			t.Fatal("expected error")
		}
	}
}

func TestFieldMissing(t *testing.T) {
	const source = memory
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeFields(db,
			[]pair{{"T", ""}},
			nil,
			"z"); e != nil {
			t.Fatal(e)
		} else {
			var missing []string
			if e := MissingFields(db, func(n string) (err error) {
				missing = append(missing, n)
				return
			}); e != nil {
				t.Fatal(e)
			}
			if !reflect.DeepEqual(missing, []string{"z"}) {
				t.Fatal("expected match", missing)
			} else {
				t.Log("okay, missing", missing)
			}
		}
	}
}
