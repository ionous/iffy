package assembly

import (
	"database/sql"
	"os/user"
	"path"
	"reflect"
	"strconv"
	"testing"

	"github.com/ionous/iffy/ephemera"
	_ "github.com/mattn/go-sqlite3"
)

func getPath(file string) (ret string, err error) {
	if user, e := user.Current(); e != nil {
		err = e
	} else {
		ret = path.Join(user.HomeDir, file)
	}
	return
}

// TestAncestors verifies valid parent-child ephemera can generate a valid ancestry table.
func TestAncestors(t *testing.T) {
	if db, e := sql.Open("sqlite3", memory); e != nil {
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
	if db, e := sql.Open("sqlite3", memory); e != nil {
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
	if db, e := sql.Open("sqlite3", memory); e != nil {
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

// TestMissingKinds to verify the kinds mentioned in parent-child ephemera exist.
func TestMissingKinds(t *testing.T) {
	if db, e := sql.Open("sqlite3", memory); e != nil {
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
		//
		w := NewModeler(dbq)
		for k, v := range kinds.cache {
			k, path := k, v.GetAncestors()
			if e := w.WriteAncestor(k, path); e != nil {
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
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		rec := ephemera.NewRecorder("ancestorTest", dbq)

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
