package assembly

import (
	"strconv"
	"testing"

	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

// TestAncestors verifies valid parent-child ephemera can generate a valid ancestry table.
func TestAncestors(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
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
		addKinds(asm.rec, pairs)
		var kinds cachedKinds
		if e := kinds.AddAncestorsOf(asm.db, "T"); e != nil {
			t.Fatal(e)
		}
		for name, kind := range kinds.cache {
			t.Log(name, ":", kind.GetAncestors())
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
		for name, v := range map[string]string{
			// kind, ancestors
			"T": "",
			"Q": "T",
			"J": "Q,T",
			"P": "J,Q,T",
			"K": "P,J,Q,T",
			"L": "P,J,Q,T",
			"M": "L,P,J,Q,T",
		} {
			kind := kinds.Get(name)
			if a := kind.GetAncestors(); a != v {
				t.Fatal("expected", v, "have", a)
			}
		}
	}
}

// TestAncestorCycle verifies cycles in parent-child ephemera generate errors.
// ex. P inherits from T; T inherits from P.
func TestAncestorCycle(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		addKinds(asm.rec, []string{
			// kind, ancestor
			"P", "T",
			"T", "P",
		})
		//
		var kinds cachedKinds
		if e := kinds.AddAncestorsOf(asm.db, "T"); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

// TestAncestorConflict verifies conflicting parent ephemera (multiple inheritance) generates an error.
// ex. P,Q inherits from T; K inherits from P and Q.
func TestAncestorConflict(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		addKinds(asm.rec, []string{
			// kind, ancestor
			"P", "T",
			"Q", "T",
			"K", "P",
			"K", "Q",
		})
		var kinds cachedKinds
		if e := kinds.AddAncestorsOf(asm.db, "T"); e == nil {
			for name, kind := range kinds.cache {
				t.Log(name, ":", kind.GetAncestors())
			}
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

func addKinds(rec *ephemera.Recorder, pairs []string) {
	for i := 0; i < len(pairs); i += 2 {
		kid := rec.NewName(pairs[i], tables.NAMED_KINDS, strconv.Itoa(i))
		ancestor := rec.NewName(pairs[i+1], tables.NAMED_KINDS, strconv.Itoa(i+1))
		rec.NewKind(kid, ancestor)
	}
}
