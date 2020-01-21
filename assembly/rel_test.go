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

type dbrel struct {
	r, k, q, c string
}

// write some relation ephemera to the database
// ( from which assembly will determine relations )
func addRelations(rec *ephemera.Recorder, els []dbrel) (err error) {
	for _, el := range els {
		r := rec.Named(ephemera.NAMED_RELATION, el.r, "test")
		k := rec.Named(ephemera.NAMED_KIND, el.k, "test")
		q := rec.Named(ephemera.NAMED_KIND, el.q, "test")
		c := el.c
		rec.NewRelation(r, k, q, c)
	}
	return
}

func matchRelations(db *sql.DB, want []dbrel) (err error) {
	var curr dbrel
	var have []dbrel
	if e := dbutil.QueryAll(db,
		`select relation, kind, otherKind, cardinality
			from mdl_rel 
			order by relation, kind, otherKind, cardinality`,
		func() (err error) {
			have = append(have, curr)
			return
		}, &curr.r, &curr.k, &curr.q, &curr.c); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch", "have:", pretty.Sprint(have), "want:", pretty.Sprint(want))
	}
	return
}

// TestRelationCreation to verify it's possible to build relations
func TestRelationCreation(t *testing.T) {
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		w := NewModeler(dbq)
		if e := fakeHierarchy(w, []pair{
			{"T", ""},
			{"P", "T"},
			{"Q", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addRelations(
			ephemera.NewRecorder("TestRelationCreation", dbq), []dbrel{
				{"R", "P", "Q", ephemera.ONE_TO_MANY},
				{"G", "P", "Q", ephemera.MANY_TO_ONE},
				{"G", "P", "Q", ephemera.MANY_TO_ONE},
				{"H", "P", "P", ephemera.ONE_TO_MANY},
			}); e != nil {
			t.Fatal(e)
		} else if e := DetermineRelations(w, db); e != nil {
			t.Fatal(e)
		} else if e := matchRelations(db, []dbrel{
			{"G", "P", "Q", ephemera.MANY_TO_ONE},
			{"H", "P", "P", ephemera.ONE_TO_MANY},
			{"R", "P", "Q", ephemera.ONE_TO_MANY},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestRelationCardinality detects conflicting cardinalities
func TestRelationCardinality(t *testing.T) {
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		w := NewModeler(dbq)
		if e := fakeHierarchy(w, []pair{
			{"P", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := addRelations(
			ephemera.NewRecorder("TestRelationCreation", dbq), []dbrel{
				{"R", "P", "P", ephemera.ONE_TO_MANY},
				{"R", "P", "P", ephemera.MANY_TO_ONE},
			}); e != nil {
			t.Fatal(e)
		} else if e := DetermineRelations(w, db); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

// TestRelationLca
func TestRelationLcaSuccess(t *testing.T) {
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		w := NewModeler(dbq)
		if e := fakeHierarchy(w, []pair{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
			{"D", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addRelations(
			ephemera.NewRecorder("TestRelationCreation", dbq), []dbrel{
				{"R", "P", "T", ephemera.ONE_TO_MANY},
				{"R", "D", "T", ephemera.ONE_TO_MANY},
				{"R", "C", "T", ephemera.ONE_TO_MANY},
			}); e != nil {
			t.Fatal(e)
		} else if e := DetermineRelations(w, db); e != nil {
			t.Fatal(e)
		} else if e := matchRelations(db, []dbrel{
			{"R", "P", "T", ephemera.ONE_TO_MANY},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestRelationLca
func TestRelationLcaFail(t *testing.T) {
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		w := NewModeler(dbq)
		if e := fakeHierarchy(w, []pair{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
			{"D", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addRelations(
			ephemera.NewRecorder("TestRelationCreation", dbq), []dbrel{
				{"R", "D", "T", ephemera.ONE_TO_MANY},
				{"R", "C", "T", ephemera.ONE_TO_MANY},
			}); e != nil {
			t.Fatal(e)
		} else if e := DetermineRelations(w, db); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}
