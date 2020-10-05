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
	t.Run("success", func(t *testing.T) {
		if asm, e := newAssemblyTest(t, memory); e != nil {
			t.Fatal(e)
		} else {
			defer asm.db.Close()
			if e := AddTestHierarchy(asm.assembler,
				"Ts", "",
			); e != nil {
				t.Fatal(e)
			} else if e := addNounEphemera(asm.rec,
				"apple", "T",
				"pear", "T",
				"toy boat", "T",
			); e != nil {
				t.Fatal(e)
			} else if e := AssembleNouns(asm.assembler); e != nil {
				t.Fatal(e)
			} else if e := matchNouns(asm.db, []modeledNoun{
				{"apple", "Ts", 0},
				{"pear", "Ts", 0},
				{"toy boat", "Ts", 0},
				{"toyBoat", "Ts", 1},
				{"boat", "Ts", 2},
				{"toy", "Ts", 3},
			}); e != nil {
				t.Fatal(e)
			}
		}
	})
	t.Run("failure", func(t *testing.T) {
		if asm, e := newAssemblyTest(t, memory); e != nil {
			t.Fatal(e)
		} else {
			defer asm.db.Close()
			if e := AddTestHierarchy(asm.assembler,
				"Ts", "",
			); e != nil {
				t.Fatal(e)
			} else if e := addNounEphemera(asm.rec,
				"bad apple", "Bs",
			); e != nil {
				t.Fatal(e)
			} else if e := AssembleNouns(asm.assembler); e == nil {
				t.Fatal("expected error")
			} else if !containsOnly(asm.dilemmas, `missing noun: "bad apple"`) {
				t.Fatal(asm.dilemmas)
			}
		}
	})
}

func collectNouns(db *sql.DB) (ret []modeledNoun, err error) {
	var curr modeledNoun
	var nouns []modeledNoun
	if e := tables.QueryAll(db,
		`select me.name, mn.kind, me.rank
		from mdl_name me
		join mdl_noun mn
			using (noun)
		order by me.noun, me.rank, me.name`,
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

func addNounEphemera(rec *ephemera.Recorder, els ...string) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		key, value := els[i], els[i+1]
		n := rec.NewName(key, tables.NAMED_NOUN, "test")
		k := rec.NewName(value, tables.NAMED_KIND, "test")
		rec.NewNoun(n, k)
	}
	return
}

// TestNounLcaSuccess to verify we can successfully determine the lowest common ancestor of nouns.
func TestNounLcaSuccess(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
			"Ps", "Ts",
			"Cs", "Ps,Ts",
			"Ds", "Ps,Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := addNounEphemera(asm.rec,
			"apple", "C",
			"apple", "P",
			"pear", "D",
			"pear", "T",
			"bandanna", "C",
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleNouns(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchNouns(asm.db, []modeledNoun{
			{"apple", "Ps", 0},
			{"bandanna", "Cs", 0},
			{"pear", "Ts", 0},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestNounLcaFailure to verify a mismatched noun hierarchy generates an error.
func TestNounLcaFailure(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
			"Ps", "Ts",
			"Cs", "Ps,Ts",
			"Ds", "Ps,Ts",
		); e != nil {
			t.Fatal(e)
		} else if e := addNounEphemera(asm.rec,
			"apple", "C",
			"apple", "D",
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleNouns(asm.assembler); e == nil {
			t.Fatal("expected failure")
		} else {
			t.Log("okay:", e)
		}
	}
}

// TestNounParts to verify a single noun generates multi part names
func TestNounParts(t *testing.T) {
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.assembler,
			"Ts", "",
		); e != nil {
			t.Fatal(e)
		} else if e := addNounEphemera(asm.rec,
			"collection of words", "T",
		); e != nil {
			t.Fatal(e)
		} else if e := AssembleNouns(asm.assembler); e != nil {
			t.Fatal(e)
		} else if e := matchNouns(asm.db, []modeledNoun{
			{"collection of words", "Ts", 0},
			{"collectionOfWords", "Ts", 1},
			{"words", "Ts", 2},
			{"of", "Ts", 3},
			{"collection", "Ts", 4},
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
