package assembly

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

// TestAncestors verifies valid parent-child ephemera can generate a valid ancestry cache.
// these arent necessarily valid or logical ancestries
func TestAncestors(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		kinds := []string{
			// kid, ancestor
			"T", "",
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
		addKinds(asm, kinds)
		t.Run("csv", func(t *testing.T) {
			var buf strings.Builder
			tables.WriteCsv(asm.db, &buf, `select * from asm_ancestry order by parent, kid`, 2)
			if have, want := buf.String(), lines(
				// parent, kid
				"Js,Ms",
				"Js,Ps",
				"Ls,Ms",
				"Ps,Ks",
				"Ps,Ls",
				"Qs,Js",
				"Qs,Ks",
				"Ts,Ps",
				"Ts,Qs",
			); have != want {
				t.Fatal(have)
			}
		})
		t.Run("kidsOf", func(t *testing.T) {
			// test raw kids of various parents
			want := []string{
				// parent, kids
				"Js", "Ms,Ps",
				"Ks", "",
				"Ls", "Ms",
				"Ms", "",
				"Ps", "Ks,Ls",
				"Qs", "Js,Ks",
				"Ts", "Ps,Qs",
			}
			for i := 0; i < len(want); i += 2 {
				parent, expected := want[i], want[i+1]
				var kids []string
				if e := ephemera.KidsOf(asm.db, parent, func(k string) {
					kids = append(kids, k)
				}); e != nil {
					t.Fatal(e)
				}
				sort.Strings(kids)
				if have := strings.Join(kids, ","); have != expected {
					t.Fatalf("%q shouldnt have kids %q", parent, have)
				}
			}
		})
		t.Run("cache", func(t *testing.T) {
			var cache cachedKinds
			if e := cache.AddDescendentsOf(asm.db, "Ts"); e != nil {
				t.Fatal(e)
			}
			// verify our original expectations: but not the root Ts slot
			for i := 2; i < len(kinds); i += 2 {
				kid := cache.Get(kinds[i] + "s")
				ancestor := cache.Get(kinds[i+1] + "s")
				if !kid.HasAncestor(ancestor) {
					t.Fatal(ancestor, "should be an ancestor of", kid)
				}
			}
			// verify our expected tree
			for name, v := range map[string]string{
				// kid, ancestors
				"Ts": "",
				"Qs": "Ts",
				"Js": "Qs,Ts",
				"Ps": "Js,Qs,Ts",
				"Ks": "Ps,Js,Qs,Ts",
				"Ls": "Ps,Js,Qs,Ts",
				"Ms": "Ls,Ps,Js,Qs,Ts",
			} {
				if a := cache.Get(name).GetAncestors(); a != v {
					t.Fatal("expected", v, "have", a)
				}
			}
		})
	}
}

// TestAncestorCycle verifies cycles in parent-child ephemera generate errors.
// ex. P inherits from T; T inherits from P.
func TestAncestorCycle(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		addKinds(asm, []string{
			// kid, ancestor
			"T", "",
			"P", "T",
			"T", "P",
		})
		//
		var kinds cachedKinds
		if e := kinds.AddDescendentsOf(asm.db, "Ts"); e == nil {
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
		addKinds(asm, []string{
			// kid, ancestor
			"T", "",
			"P", "T",
			"Q", "T",
			"K", "P",
			"K", "Q",
		})
		var kinds cachedKinds
		if e := kinds.AddDescendentsOf(asm.db, "Ts"); e == nil {
			t.Fatal("expected error")
		} else {
			t.Log("okay:", e)
		}
	}
}

func addKinds(asm *assemblyTest, pairs []string) {
	plural := make(map[string]bool)
	rec, w := asm.rec, asm.assembler
	for _, n := range pairs {
		// normally, we would assemble plurals first
		// pluralizing singular names
		if len(n) > 0 && !plural[n] {
			plural[n] = true
			// one, many
			w.WritePlural(n, n+"s")
		}
	}

	for i := 0; i < len(pairs); i += 2 {
		// cats are a kind of animal
		pluralKid, ancestor := pairs[i]+"s", pairs[i+1]
		if len(ancestor) > 0 {
			kr := rec.NewName(pluralKid, tables.NAMED_KINDS, strconv.Itoa(i))
			ar := rec.NewName(ancestor, tables.NAMED_KIND, strconv.Itoa(i+1))
			rec.NewKind(kr, ar)
		}
	}
}
