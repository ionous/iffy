package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// TestNounFormation to verify we can successfully assemble nouns from ephemera
func TestNounFormation(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		//
		if e := AddTestHierarchy(m, []TargetField{
			{"T", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := addNouns(rec, []pair{
			{"apple", "T"},
			{"pear", "T"},
			{"toy boat", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineNouns(m, db); e != nil {
			t.Fatal(e)
		} else if e := matchNouns(db, []modeledNoun{
			{"apple", "T", 0},
			{"pear", "T", 0},
			{"toy boat", "T", 0},
			{"boat", "T", 1},
			{"toy", "T", 2},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

func collectNouns(db *sql.DB) (ret []modeledNoun, err error) {
	var curr modeledNoun
	var nouns []modeledNoun
	if e := tables.QueryAll(db,
		`select n.name, i.kind, n.rank
		from mdl_name n 
		join mdl_noun i
			on (n.noun = i.noun)
		order by i.noun, rank, name`,
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
		n := rec.Named(tables.NAMED_NOUN, el.key, "test")
		k := rec.Named(tables.NAMED_KIND, el.value, "test")
		rec.NewNoun(n, k)
	}
	return
}

// TestNounLcaSucess to verify we can successfully determine the lowest common ancestor of nouns.
func TestNounLcaSucess(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		//
		if e := AddTestHierarchy(m, []TargetField{
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
		} else if e := DetermineNouns(m, db); e != nil {
			t.Fatal(e)
		} else if e := matchNouns(db, []modeledNoun{
			{"apple", "P", 0},
			{"bandanna", "C", 0},
			{"pear", "T", 0},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestNounLcaFailure to verify a mismatched noun hierarchy generates an error.
func TestNounLcaFailure(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		//
		if e := AddTestHierarchy(m, []TargetField{
			{"T", ""},
			{"P", "T"},
			{"C", "P,T"},
			{"D", "P,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addNouns(rec, []pair{
			{"apple", "C"},
			{"apple", "D"},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineNouns(m, db); e == nil {
			t.Fatal("expected failure")
		} else {
			t.Log("okay:", e)
		}
	}
}

// TestNounParts to verify a single noun generates multi part names
func TestNounParts(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		db, rec, m := t.db, t.rec, t.modeler
		//
		if e := AddTestHierarchy(m, []TargetField{
			{"T", ""},
		}); e != nil {
			t.Fatal(e)
		} else if e := addNouns(rec, []pair{
			{"collection of words", "T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineNouns(m, db); e != nil {
			t.Fatal(e)
		} else if e := matchNouns(db, []modeledNoun{
			{"collection of words", "T", 0},
			{"words", "T", 1},
			{"of", "T", 2},
			{"collection", "T", 3},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

func matchNouns(db *sql.DB, want []modeledNoun) (err error) {
	if got, e := collectNouns(db); e != nil {
		err = e
	} else if !reflect.DeepEqual(got, want) {
		e := errutil.New("mismatch",
			"have:", pretty.Sprint(got),
			"want:", pretty.Sprint(want))
		err = e
	}
	return
}
