package assembly

import (
	"strconv"
	"testing"

	"github.com/ionous/iffy/ephemera"
	_ "github.com/mattn/go-sqlite3"
)

// TestAncestors verifies valid parent-child ephemera can generate a valid ancestry table.
func TestAncestors(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec := t.db, t.rec
		//
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
			kid := rec.Named(ephemera.NAMED_KIND, pairs[i], strconv.Itoa(i))
			ancestor := rec.Named(ephemera.NAMED_KIND, pairs[i+1], strconv.Itoa(i+1))
			rec.NewKind(kid, ancestor)
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

// TestAncestorCycle verifies cycles in parent-child ephemera generate errors.
// ex. P inherits from T; T inherits from P.
func TestAncestorCycle(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec := t.db, t.rec
		//
		pairs := []string{
			// kind, ancestor
			"P", "T",
			"T", "P",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := rec.Named(ephemera.NAMED_KIND, pairs[i], strconv.Itoa(i))
			parent := rec.Named(ephemera.NAMED_KIND, pairs[i+1], strconv.Itoa(i+1))
			rec.NewKind(kid, parent)
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

// TestAncestorConflict verifies conflicting parent ephemera (multiple inheritance) generates an error.
// ex. P,Q inherits from T; K inherits from P and Q.
func TestAncestorConflict(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec := t.db, t.rec
		//
		pairs := []string{
			// kind, ancestor
			"P", "T",
			"Q", "T",
			"K", "P",
			"K", "Q",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := rec.Named(ephemera.NAMED_KIND, pairs[i], strconv.Itoa(i))
			parent := rec.Named(ephemera.NAMED_KIND, pairs[i+1], strconv.Itoa(i+1))
			rec.NewKind(kid, parent)
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
