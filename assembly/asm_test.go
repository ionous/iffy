package assembly

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/ionous/iffy/ephemera"
	_ "github.com/mattn/go-sqlite3"
)

func TestAncestors(t *testing.T) {
	const source = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		rec := ephemera.NewRecorder("ancestorTest", dbq)
		pairs := []string{
			// kind, ancestor
			"P", "T",
			"Q", "T",
			"L", "P",
			"K", "P",
			"K", "Q",
			"J", "Q",
			"M", "L",
			"P", "J",
			"M", "J",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := rec.Named("kind", pairs[i], strconv.Itoa(i))
			ancestor := rec.Named("kind", pairs[i+1], strconv.Itoa(i+1))
			rec.Kind(kid, ancestor)
		}
		//
		kinds := &cachedKinds{}
		if e := kinds.AddAncestorsOf(db, "T"); e != nil {
			t.Fatal(e)
		}
		for k, n := range kinds.cache {
			t.Log(k, ":", n.GetAncestors())
		}
		// verify our original expectations
		for i := 0; i < len(pairs); i += 2 {
			kid := kinds.Get(pairs[i])
			ancestor := kinds.Get(pairs[i+1])
			if !kid.HasAncestor(ancestor) {
				t.Fatal(ancestor, "should be an ancestor of", kid)
			}
		}
		// verify our expected tree
		for k, v := range map[string]string{
			// kind, ancestors
			"T": "",
			"Q": "T",
			"J": "Q,T",
			"P": "J,Q,T",
			"K": "P,J,Q,T",
			"L": "P,J,Q,T",
			"M": "L,P,J,Q,T",
		} {
			k := kinds.Get(k)
			if a := k.GetAncestors(); a != v {
				t.Fatal("expected", v, "have", a)
			}
		}
	}
}

func TestAncestorCycle(t *testing.T) {
	const source = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		rec := ephemera.NewRecorder("ancestorTest", dbq)
		pairs := []string{
			// kind, ancestor
			"P", "T",
			"T", "P",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := rec.Named("kind", pairs[i], strconv.Itoa(i))
			parent := rec.Named("kind", pairs[i+1], strconv.Itoa(i+1))
			rec.Kind(kid, parent)
		}
		//
		kinds := &cachedKinds{}
		if e := kinds.AddAncestorsOf(db, "T"); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

func TestAncestorConflict(t *testing.T) {
	const source = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		rec := ephemera.NewRecorder("ancestorTest", dbq)
		pairs := []string{
			// kind, ancestor
			"P", "T",
			"Q", "T",
			"K", "P",
			"K", "Q",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := rec.Named("kind", pairs[i], strconv.Itoa(i))
			parent := rec.Named("kind", pairs[i+1], strconv.Itoa(i+1))
			rec.Kind(kid, parent)
		}
		//
		kinds := &cachedKinds{}
		if e := kinds.AddAncestorsOf(db, "T"); e == nil {
			for k, n := range kinds.cache {
				t.Log(k, ":", n.GetAncestors())
			}
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

func TestMissingKinds(t *testing.T) {
	const source = "file:test.db?cache=shared&mode=memory"
	/*if source, e := getPath(); e != nil {
		t.Fatal(e)
	} else*/if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		rec := ephemera.NewRecorder("ancestorTest", dbq)
		pairs := []string{
			// kind, ancestor
			"P", "T",
			"Q", "T",
			"P", "R",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := rec.Named("kind", pairs[i], strconv.Itoa(i))
			parent := rec.Named("kind", pairs[i+1], strconv.Itoa(i+1))
			rec.Kind(kid, parent)
		}
		// add the kinds
		kinds := &cachedKinds{}
		if e := kinds.AddAncestorsOf(db, "T"); e != nil {
			for k, n := range kinds.cache {
				t.Log(k, ":", n.GetAncestors())
			}
			t.Fatal(e)
		}
		//
		w := NewWriter(dbq)
		for k, v := range kinds.cache {
			k, path := k, v.GetAncestors()
			w.WriteAncestor(k, path)
		}
		// now test for our missing "R"
		var missing []string
		if e := MissingKinds(db, func(k string) {
			missing = append(missing, k)
		}); e != nil {
			t.Fatal(e)
		}
		if len(missing) != 1 || missing[0] != "R" {
			t.Fatal("expected R, have", missing)
		}

	}
}
