package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/kr/pretty"
)

// TestNounFormation to verify we can successfully assemble nouns from ephemera
func TestNounFormation(t *testing.T) {
	const source = memory
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		rec := ephemera.NewRecorder(t.Name(), dbq)
		w := NewModeler(dbq)
		if e := fakeHierarchy(w, []pair{
			{"T", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := addNouns(rec, []pair{
			{"apple", "T"},
			{"pear", "T"},
			{"machine gun", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineNouns(w, db); e != nil {
			t.Fatal(e)
		} else if got, e := collectNouns(db); e != nil {
			t.Fatal(e)
		} else {
			want := []modeledNoun{
				{"apple", "T", 0},
				{"gun", "T", 1},
				{"machine", "T", 2},
				{"machine gun", "T", 0},
				{"pear", "T", 0},
			}
			if !reflect.DeepEqual(got, want) {
				e := errutil.New("mismatch",
					"have:", pretty.Sprint(got),
					"want:", pretty.Sprint(want))
				t.Fatal(e)
			}
		}
	}
}

func collectNouns(db *sql.DB) (ret []modeledNoun, err error) {
	var curr modeledNoun
	var nouns []modeledNoun
	if e := dbutil.QueryAll(db,
		`select n.name, i.kind, n.rank
		from mdl_name n join mdl_noun i
			on (n.idModelNoun = i.id)
		order by name`,
		func() (err error) {
			nouns = append(nouns, curr)
			return
		}, &curr.name, &curr.kind, &curr.rank); e != nil {
		err = e
	} else {
		ret = nouns
	}
	return
}

type modeledNoun struct {
	name, kind string
	rank       int
}

func addNouns(rec *ephemera.Recorder, els []pair) (err error) {
	for _, el := range els {
		n := rec.Named(ephemera.NAMED_NOUN, el.key, "test")
		k := rec.Named(ephemera.NAMED_KIND, el.value, "test")
		rec.NewNoun(n, k)
	}
	return
}

// TestNounLca to verify we can successfully determine the lowest common ancestor of nouns.
func TestNounLca(t *testing.T) {
	const source = memory
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		defer db.Close()
		dbq := ephemera.NewDBQueue(db)
		rec := ephemera.NewRecorder(t.Name(), dbq)
		w := NewModeler(dbq)
		if e := fakeHierarchy(w, []pair{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
			{"D", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addNouns(rec, []pair{
			{"apple", "C"},
			{"apple", "P"},
			{"pear", "D"},
			{"pear", "T"},
			{"bandanna", "C"},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineNouns(w, db); e != nil {
			t.Fatal(e)
		} else if got, e := collectNouns(db); e != nil {
			t.Fatal(e)
		} else {
			want := []modeledNoun{
				{"apple", "P", 0},
				{"bandanna", "C", 0},
				{"pear", "T", 0},
			}
			if !reflect.DeepEqual(got, want) {
				e := errutil.New("mismatch",
					"have:", pretty.Sprint(got),
					"want:", pretty.Sprint(want))
				t.Fatal(e)
			}
		}
	}
}

// noun failed lca
// multipart noun name
