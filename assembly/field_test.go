package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
)

func TestLca(t *testing.T) {
	match := func(a, b, c []string) bool {
		chain := findOverlap(a, b)
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
type kp struct{ kind, path string }

func writeFields(db *sql.DB, kinds []kp, kfps []kfp, missing ...string) (err error) {
	dbq := ephemera.NewDBQueue(db)
	r := ephemera.NewRecorder("ancestorTest", dbq)
	w := NewWriter(dbq)
	// create some hierarchy
	for _, p := range kinds {
		w.WriteAncestor(p.kind, p.path)
	}
	// write some primitives
	for _, p := range kfps {
		kind := r.Named(ephemera.NAMED_KIND, p.kind, "test")
		field := r.Named(ephemera.NAMED_FIELD, p.field, "test")
		r.Primitive(p.fieldType, kind, field)
	}
	for _, m := range missing {
		r.Named(ephemera.NAMED_FIELD, m, "test")
	}
	if e := DetermineFields(w, db); e != nil {
		err = e
	}
	return
}
func matchProperties(db *sql.DB, expected []kfp) (err error) {
	if it, e := db.Query(`select kind,field,type from property order by kind, field, type`); e != nil {
		err = e
	} else {
		defer it.Close()
		for cnt := 0; it.Next(); cnt++ {
			var curr kfp
			if e := it.Scan(&curr.kind, &curr.field, &curr.fieldType); e != nil {
				err = e
				break
			} else {
				want := expected[cnt]
				if !reflect.DeepEqual(curr, want) {
					err = errutil.New("mismatch", "have:", pretty.Sprint(curr), "want:", pretty.Sprint(want))
					break
				}
			}
		}
		// tests if early exit
		if e := it.Err(); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func TestFields(t *testing.T) {
	const source = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeFields(db,
			[]kp{
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
	const source = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeFields(db,
			[]kp{
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

func TestFieldTypeMismatch(t *testing.T) {
	const source = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeFields(db,
			[]kp{
				{"T", ""},
			},
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

// FIX FIX missing properties ( named but not specified )
// select from named where kind == NAMED_FIELD not exists in
func TestFieldMissing(t *testing.T) {
	const source = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		if e := writeFields(db,
			[]kp{{"T", ""}},
			nil,
			"z"); e != nil {
			t.Fatal(e)
		} else {
			var missing []string
			if e := MissingFields(db, func(n string) {
				missing = append(missing, n)
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
