package assembly

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/ionous/iffy/tables"
)

// TestMissingKinds to verify the kinds mentioned in parent-child ephemera exist.
func TestMissingKinds(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		db, rec, m := asm.db, asm.rec, asm.modeler
		// kind, ancestor
		pairs := []string{
			"P", "T",
			"Q", "T",
			"P", "R",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := rec.NewName(pairs[i], tables.NAMED_KIND, strconv.Itoa(i))
			parent := rec.NewName(pairs[i+1], tables.NAMED_KIND, strconv.Itoa(i+1))
			rec.NewKind(kid, parent)
		}
		// add the kinds
		kinds := &cachedKinds{}
		if e := kinds.AddAncestorsOf(db, "T"); e != nil {
			for k, n := range kinds.cache {
				t.Log(k, ":", n.GetAncestors())
			}
			t.Fatal(e)
		}
		for k, v := range kinds.cache {
			k, path := k, v.GetAncestors()
			if e := m.WriteAncestor(k, path); e != nil {
				t.Fatal(e)
			}
		}
		// now test for our missing "R"
		var missing []string
		if e := MissingKinds(db, func(k string) (err error) {
			missing = append(missing, k)
			return
		}); e != nil {
			t.Fatal(e)
		}
		if len(missing) != 1 || missing[0] != "R" {
			t.Fatal("expected R, have", missing)
		}
	}
}

// TestMissingAspects detects fields labeled as aspects which are missing from the aspects ephemera.
func TestMissingAspects(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		db, rec := asm.db, asm.rec
		//
		parent := rec.NewName("K", tables.NAMED_KIND, "container")
		for i, aspect := range []string{
			// known, unknown
			"A", "F",
			"C", "D",
			"E", "B",
		} {
			a := rec.NewName(aspect, tables.NAMED_ASPECT, "test")
			if known := i&1 == 0; known {
				rec.NewAspect(a)
			}
			rec.NewPrimitive(tables.PRIM_ASPECT, parent, a)
		}
		expected := []string{"B", "D", "F"}
		if missing, e := undeclaredAspects(db); e != nil {
			t.Fatal(e)
		} else if matches := reflect.DeepEqual(missing, expected); !matches {
			t.Fatal("want:", expected, "have:", missing)
		} else {
			t.Log("okay")
		}
	}
}

func TestMissingField(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.modeler, []TargetField{
			{"T", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeMissing(asm.rec, []string{
			"z",
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineFields(asm.modeler, asm.db); e != nil {
			t.Fatal(e)
		} else {
			var missing []string
			if e := MissingFields(asm.db, func(n string) (err error) {
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

// xTestMissingUnknownField missing properties ( kind, field pair doesn't exist in model )
func xTestMissingUnknownField(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.modeler, []TargetField{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestFields(asm.modeler, []TargetValue{
			{"T", "d", tables.PRIM_DIGI},
			{"T", "t", tables.PRIM_TEXT},
			{"T", "t2", tables.PRIM_TEXT},
			{"P", "p", tables.PRIM_TEXT},
			{"C", "c", tables.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := addDefaults(asm.rec, []triplet{
			{"T", "t", "some text"},
			{"P", "t2", "other text"},
			{"C", "c", "c text"},
			{"P", "c", "invalid"}, // this pair doesnt exist
			{"T", "p", "invalid"}, // this pair doesnt exist
			{"C", "d", 123},
		}); e != nil {
			t.Fatal(e)
		} else {
			var got []pair
			if e := MissingDefaults(asm.db, func(k, f string) (err error) {
				got = append(got, pair{k, f})
				return
			}); e != nil {
				t.Fatal(e)
			} else if !reflect.DeepEqual(got, []pair{
				{"P", "c"},
				{"T", "p"},
			}) {
				t.Fatal("mismatched", got)
			}
		}
	}
}
