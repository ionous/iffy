package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// TestVerbMismatches verifies that we can collapse multiple relation-verb pairs so long as the verb-stem pair match
// while ensuring the same stem cannot be used in multiple relations.
func TestVerbMismatches(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		m := t.modeler
		//
		if e := m.WriteVerb("R", "contains"); e != nil {
			t.Fatal(e)
		} else if e := m.WriteVerb("R", "containing"); e != nil {
			t.Fatal(e)
		} else if e := m.WriteVerb("Q", "supporting"); e != nil {
			t.Fatal(e)
		} else if e := m.WriteVerb("Q", "supports"); e != nil {
			t.Fatal(e)
		} else if e := m.WriteVerb("R", "supports"); e == nil {
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
		r := rec.Named(tables.NAMED_RELATION, el.r, "test")
		k := rec.Named(tables.NAMED_KIND, el.k, "test")
		q := rec.Named(tables.NAMED_KIND, el.q, "test")
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
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		//
		if e := fakeHierarchy(m, []pair{
			{"T", ""},
			{"P", "T"},
			{"Q", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addRelations(
			rec, []dbrel{
				{"R", "P", "Q", tables.ONE_TO_MANY},
				{"G", "P", "Q", tables.MANY_TO_ONE},
				{"G", "P", "Q", tables.MANY_TO_ONE},
				{"H", "P", "P", tables.ONE_TO_MANY},
			}); e != nil {
			t.Fatal(e)
		} else if e := DetermineRelations(m, db); e != nil {
			t.Fatal(e)
		} else if e := matchRelations(db, []dbrel{
			{"G", "P", "Q", tables.MANY_TO_ONE},
			{"H", "P", "P", tables.ONE_TO_MANY},
			{"R", "P", "Q", tables.ONE_TO_MANY},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestRelationCardinality detects conflicting cardinalities
func TestRelationCardinality(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		//
		if e := fakeHierarchy(m, []pair{
			{"P", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := addRelations(rec, []dbrel{
			{"R", "P", "P", tables.ONE_TO_MANY},
			{"R", "P", "P", tables.MANY_TO_ONE},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineRelations(m, db); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

// TestRelationLca
func TestRelationLcaSuccess(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		//
		if e := fakeHierarchy(m, []pair{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
			{"D", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addRelations(rec, []dbrel{
			{"R", "P", "T", tables.ONE_TO_MANY},
			{"R", "D", "T", tables.ONE_TO_MANY},
			{"R", "C", "T", tables.ONE_TO_MANY},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineRelations(m, db); e != nil {
			t.Fatal(e)
		} else if e := matchRelations(db, []dbrel{
			{"R", "P", "T", tables.ONE_TO_MANY},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestRelationLcaFailure to verify a mismatched relation hierarchy generates an error.
func TestRelationLcaFailure(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		//
		if e := fakeHierarchy(m, []pair{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
			{"D", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addRelations(rec, []dbrel{
			{"R", "D", "T", tables.ONE_TO_MANY},
			{"R", "C", "T", tables.ONE_TO_MANY},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineRelations(m, db); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}
