package assembly

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/ionous/iffy/ephemera"
	_ "github.com/mattn/go-sqlite3"
)

// TestMissingKinds to verify the kinds mentioned in parent-child ephemera exist.
func TestMissingKinds(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		// kind, ancestor
		pairs := []string{
			"P", "T",
			"Q", "T",
			"P", "R",
		}
		for i := 0; i < len(pairs); i += 2 {
			kid := rec.Named(ephemera.NAMED_KIND, pairs[i], strconv.Itoa(i))
			parent := rec.Named(ephemera.NAMED_KIND, pairs[i+1], strconv.Itoa(i+1))
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
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec := t.db, t.rec
		//
		parent := rec.Named(ephemera.NAMED_KIND, "K", "container")
		for i, aspect := range []string{
			// known, unknown
			"A", "F",
			"C", "D",
			"E", "B",
		} {
			a := rec.Named(ephemera.NAMED_ASPECT, aspect, "test")
			if known := i&1 == 0; known {
				rec.NewAspect(a)
			}
			rec.NewPrimitive(ephemera.PRIM_ASPECT, parent, a)
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
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := fakeHierarchy(t.modeler, []pair{
			{"T", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := writeMissing(t.rec, []string{
			"z",
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineFields(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else {
			var missing []string
			if e := MissingFields(t.db, func(n string) (err error) {
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
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := fakeHierarchy(t.modeler, []pair{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeFields(t.modeler, []kfp{
			{"T", "d", ephemera.PRIM_DIGI},
			{"T", "t", ephemera.PRIM_TEXT},
			{"T", "t2", ephemera.PRIM_TEXT},
			{"P", "p", ephemera.PRIM_TEXT},
			{"C", "c", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := addDefaults(t.rec, []triplet{
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
			if e := MissingDefaults(t.db, func(k, f string) (err error) {
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
