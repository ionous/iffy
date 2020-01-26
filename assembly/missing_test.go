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
		//
		pairs := []string{
			// kind, ancestor
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
		if e := writeFields(t, []pair{
			{"T", ""},
		}, nil, "z"); e != nil {
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
